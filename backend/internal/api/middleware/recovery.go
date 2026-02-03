package middleware

import (
	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

type Recovery struct {
	loggerWriter *lumberjack.Logger
}

func CreateRecoveryMiddleware(loggerWriter *lumberjack.Logger) *Recovery {
	return &Recovery{
		loggerWriter: loggerWriter,
	}
}
func (m *Recovery) Handler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		m.loggerWriter,
		func(c *gin.Context, err interface{}) {
			response.Error(c, http.StatusInternalServerError, response.Response{
				Code:    errors.ServerError,
				Data:    nil,
				Message: http.StatusText(http.StatusInternalServerError),
			})
		},
	)
}
