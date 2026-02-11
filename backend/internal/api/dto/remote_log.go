package dto

import (
	"backend/internal/service/remote_log"
)

type LogMessage struct {
	Level   remote_log.LogLevel `json:"level" binding:"required,oneof=INFO WARN ERROR"`
	Message string              `json:"message" binding:"required"`
	Source  interface{}         `json:"source" binding:"required"`
}
