package handler

import (
	"backend/internal/api/dto"
	"backend/internal/service/remote_log"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RemoteLogHandler struct {
	svc *remote_log.LogService
	log *zap.Logger
}

func CreateRemoteLogHandler(log *zap.Logger, svc *remote_log.LogService) *RemoteLogHandler {
	return &RemoteLogHandler{
		log: log,
		svc: svc,
	}
}

func (s *RemoteLogHandler) Handler(c *gin.Context) {
	var req dto.LogMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		s.log.Error("invalid request", zap.Error(err))
		response.Error(c, http.StatusBadRequest, response.Response{
			Code:    4004,
			Data:    nil,
			Message: "invalid request",
		})
		return
	}

	s.svc.Push(req.Level, req.Message, req.Source)
	response.Success(c, gin.H{
		"code":    0,
		"message": "log accepted",
	})
}
