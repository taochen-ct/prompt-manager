package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Logger struct {
	Logger *zap.Logger
}

func CreateLoggerMiddleware(logger *zap.Logger) *Logger {
	return &Logger{
		Logger: logger,
	}
}

func (m *Logger) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		m.Logger.Info(fmt.Sprintf("request duration: %v", time.Since(start)))
	}
}
