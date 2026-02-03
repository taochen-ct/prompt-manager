package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    data,
		Message: http.StatusText(http.StatusOK),
	})
}

func Error(c *gin.Context, statusCode int, response Response) {
	c.JSON(statusCode, response)
}
