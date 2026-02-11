package remote_log

import (
	"backend/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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
}

func formatSource(src interface{}) string {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Sprintf("%v", src)
	}
	return string(b)
}
