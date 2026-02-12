package remote_log

import (
	"backend/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type LogLevel string

const (
	INFO          LogLevel = "INFO"
	WARN          LogLevel = "WARN"
	ERROR         LogLevel = "ERROR"
	DefaultBuffer int      = 1024
)

type IService interface {
	Push(level LogLevel, msg, source string)
	worker()
	handle(log LogMessage)
}

type LogMessage struct {
	Level     LogLevel    `json:"level"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Source    interface{} `json:"source"`
}

type LogService struct {
	logChan chan LogMessage
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *zap.Logger
	cfg     *config.Config
}

func CreateLogService(cfg *config.Config, logger *zap.Logger) (*LogService, func()) {
	buffer := cfg.Server.Buffer
	if buffer <= 0 {
		buffer = DefaultBuffer
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &LogService{
		logChan: make(chan LogMessage, buffer),
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger,
		cfg:     cfg,
	}

	s.wg.Add(1)
	go s.worker()

	return s, func() {
		cancel()
		close(s.logChan)
		s.wg.Wait()
	}
}

func (s *LogService) Push(level LogLevel, msg string, source interface{}) {
	select {
	case s.logChan <- LogMessage{
		Level:     level,
		Message:   msg,
		Source:    source,
		Timestamp: time.Now(),
	}:
	default:
		// 防止阻塞
		s.logger.Info("log channel full, drop log")
	}
}

func (s *LogService) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			// drain 剩余日志
			for log := range s.logChan {
				s.handle(log)
			}
			return

		case log := <-s.logChan:
			s.handle(log)
		}
	}
}

func (s *LogService) handle(log LogMessage) {
	fmt.Printf("[%s] [%s] %s (%v)\n",
		log.Timestamp.Format(time.RFC3339),
		log.Level,
		log.Message,
		formatSource(log.Source),
	)

	// TODO:
	// - 写文件
	// - 写数据库
	// - 推送 Kafka / MQ
	if s.cfg.Else.ScSend.Enable {
		title := fmt.Sprintf("%s-%s", log.Level, log.Message)
		go s.scSend(title, formatSource(log.Source))
	}
}

func formatSource(src interface{}) string {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Sprintf("%v", src)
	}
	return string(b)
}

func (s *LogService) scSend(text string, desp string) {
	data := url.Values{}
	data.Set("text", text)
	data.Set("desp", desp)
	data.Set("channel", "8|9")
	data.Set("noip", "1")

	// 根据 sendkey 是否以 "sctp" 开头决定 API 的 URL
	var apiUrl string
	key := s.cfg.Else.ScSend.Key
	if strings.HasPrefix(key, "sctp") {
		// 使用正则表达式提取数字部分
		re := regexp.MustCompile(`sctp(\d+)t`)
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 {
			num := matches[1]
			apiUrl = fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, key)
		} else {
			s.logger.Error("Invalid sendkey format for sctp")
			return
		}
	} else {
		apiUrl = fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.Error("send request create error", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("send request error", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("server response error", zap.Error(err))
		return
	}

	s.logger.Info("send success", zap.String("body", string(body)))
}
