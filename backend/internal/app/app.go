package app

import (
	"backend/pkg/config"
	"backend/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	db      *sqlx.DB
	conf    *config.Config
	logger  *zap.Logger
	proxy   *ProxyServer
	httpSrv *http.Server
}

func createHttpServer(
	conf *config.Config,
	router *gin.Engine,
) *http.Server {
	return &http.Server{
		Addr:    ":" + conf.Server.Port,
		Handler: router,
	}
}

func createApp(
	db *sqlx.DB,
	conf *config.Config,
	logger *zap.Logger,
	httpSrv *http.Server,
) (*App, error) {
	if err := runMigrate(db, conf); err != nil {
		return nil, err
	}
	return &App{
		db:      db,
		conf:    conf,
		logger:  logger,
		httpSrv: httpSrv,
	}, nil
}

func (a *App) Start() error {
	// 启动 http server
	go func() {
		log.Printf("http server started")
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(fmt.Sprintf("http server started: %s", err.Error()))
			panic(err)
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) (err error) {
	log.Printf("http server has been stop")

	if err = a.httpSrv.Shutdown(ctx); err != nil {
		a.logger.Error(fmt.Sprintf("http server shutdown error: %s", err.Error()))
		return
	}
	return
}

func Run() {
	cfg := config.Load()
	loggerWriter, loggerEnt := logger.Init(cfg)
	app, cleanup, err := wireApp(cfg, loggerWriter, loggerEnt)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 统一 context
	proxyCtx, proxyCancel := context.WithCancel(context.Background())
	defer proxyCancel()

	// 启动应用
	log.Printf("start app %s ...", app.conf.Server.Port)
	if err = app.Start(); err != nil {
		panic(err)
	}

	// 启动模型代理服务
	proxy := CreateProxyServer(cfg.Proxy)
	go proxy.Start(proxyCtx)

	// 等待中断信号以优雅地关闭应用
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("shutdown app %s ...", app.conf.Server.Port)
	proxyCancel()

	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭应用
	if err = app.Stop(ctx); err != nil {
		panic(err)
	}
}
