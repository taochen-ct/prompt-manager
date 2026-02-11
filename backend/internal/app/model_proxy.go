package app

import (
	"backend/pkg/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type APIClient struct {
	client  *http.Client
	name    string
	apiKey  string
	apiBase string
}
type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type OpenAIModelList struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}

var rrCounter = make(map[string]uint64)
var rrLock sync.Mutex
var (
	// 非 stream
	nonStreamSem chan struct{}
	// stream
	streamSem chan struct{}
)

func pickClient(model string, clients []*APIClient) *APIClient {
	rrLock.Lock()
	defer rrLock.Unlock()

	idx := rrCounter[model] % uint64(len(clients))
	rrCounter[model]++

	return clients[idx]
}

func createModelClient(name string, apiBase string, apiKey string, httpTransport *http.Transport) *APIClient {
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   0,
	}
	return &APIClient{
		client:  client,
		name:    name,
		apiKey:  apiKey,
		apiBase: apiBase,
	}
}

func openAIProxyHandler(modelClientMap map[string][]*APIClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("[PROXY] sub path = %s", path)

		// 特殊处理models路由, 返回本地代理模型
		if strings.HasPrefix(path, "/models") {
			now := time.Now().Unix()

			data := make([]OpenAIModel, 0, len(modelClientMap))
			for model := range modelClientMap {
				data = append(data, OpenAIModel{
					ID:      model,
					Object:  "model",
					Created: now,
					OwnedBy: "proxy",
				})
			}

			c.JSON(200, OpenAIModelList{
				Object: "list",
				Data:   data,
			})
			return
		}

		// 读取 body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 解析model, stream
		var payload struct {
			Model  string `json:"model"`
			Stream bool   `json:"stream"`
		}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil || payload.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing model in request body",
			})
			log.Printf("[ERROR] missing model in request body")
			return
		}

		// 限流
		var sem chan struct{}
		if payload.Stream {
			sem = streamSem
		} else {
			sem = nonStreamSem
		}
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
		default:
			c.JSON(429, gin.H{"error": "too many requests"})
			return
		}

		clients, ok := modelClientMap[payload.Model]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "unknown model: " + payload.Model,
			})
			log.Printf("[ERROR] unknown model: %s", payload.Model)
			return
		}
		// 轮询获取客户端
		client := pickClient(payload.Model, clients)
		targetURL := client.apiBase + path
		log.Printf(
			"[PROXY] client: %s target: %s; model name = %s; stream = %v",
			client.name,
			targetURL,
			payload.Model,
			payload.Stream,
		)

		// 构造转发请求
		req, err := http.NewRequest(
			c.Request.Method,
			targetURL,
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("[ERROR] %s", err.Error())
			return
		}

		// 拷贝 headers
		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		// 注入 key
		req.Header.Set("Authorization", "Bearer "+client.apiKey)

		// 转发
		resp, err := client.client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			log.Printf("[ERROR] %s", err.Error())
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for k, v := range resp.Header {
			c.Writer.Header()[k] = v
		}
		c.Writer.WriteHeader(resp.StatusCode)

		if payload.Stream {
			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")

			flusher, ok := c.Writer.(http.Flusher)
			if !ok {
				c.JSON(500, gin.H{"error": "stream not supported"})
				return
			}

			buf := make([]byte, 4096)
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					if _, werr := c.Writer.Write(buf[:n]); werr != nil {
						// client 断开
						return
					}
					flusher.Flush()
				}
				if err != nil {
					return
				}
			}
		}

		io.Copy(c.Writer, resp.Body)
	}
}

func InitModelClientMap(cfg config.Proxy, modelClientMap map[string][]*APIClient) {

	existClients := make(map[string]*APIClient)

	for modelName, mc := range cfg.Models {
		if len(mc.Endpoints) == 0 {
			log.Printf("[warn] model=%s has no endpoints, skipped", modelName)
			continue
		}
		log.Printf("[info] model=%s init", modelName)

		clients := make([]*APIClient, 0, len(mc.Endpoints))

		for idx, ep := range mc.Endpoints {
			if ep.ApiBase == "" || ep.ApiKey == "" {
				log.Printf(
					"[warn] model=%s endpoint[%d] invalid: %+v",
					modelName, idx, ep,
				)
				continue
			}

			client, ok := existClients[ep.Name]
			if !ok {
				transport := &http.Transport{
					MaxConnsPerHost:     cfg.HttpClient.MaxConnsPerHost,
					MaxIdleConns:        cfg.HttpClient.MaxIdleConns,
					MaxIdleConnsPerHost: cfg.HttpClient.MaxIdleConnsPerHost,
					IdleConnTimeout:     cfg.HttpClient.IdleConnTimeout,
				}
				client = createModelClient(ep.Name, ep.ApiBase, ep.ApiKey, transport)
				existClients[ep.Name] = client
			}

			clients = append(clients, client)

			log.Printf(
				"[init] model=%s endpoint[%d] base=%s",
				modelName, idx, ep.ApiBase,
			)
		}

		if len(clients) > 0 {
			modelClientMap[modelName] = clients
		}
	}

	log.Printf(
		"[init] modelClientMap initialized, models=%d",
		len(modelClientMap),
	)
}

func InitSemaphore(cfg config.Proxy) {
	nonStream := cfg.Limits.NonStreamConcurrency
	stream := cfg.Limits.StreamConcurrency

	// 兜底值，防止配置写 0 把服务搞挂
	if nonStream <= 0 {
		nonStream = 300
	}
	if stream <= 0 {
		stream = 50
	}

	nonStreamSem = make(chan struct{}, nonStream)
	streamSem = make(chan struct{}, stream)

	log.Printf(
		"[init] semaphore initialized nonStream=%d stream=%d",
		nonStream, stream,
	)
}

type ProxyServer struct {
	srv *http.Server
}

func CreateProxyServer(cfg config.Proxy) *ProxyServer {
	modelClientMap := make(map[string][]*APIClient, len(cfg.Models))

	InitSemaphore(cfg)
	InitModelClientMap(cfg, modelClientMap)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Any("/v1/*path", openAIProxyHandler(modelClientMap))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &ProxyServer{srv: srv}
}

func (p *ProxyServer) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Println("[proxy] shutting down...")
		_ = p.srv.Shutdown(context.Background())
	}()

	log.Printf("[proxy] listen on %s", p.srv.Addr)

	if err := p.srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Printf("[proxy] listen error: %v", err)
	}
}
