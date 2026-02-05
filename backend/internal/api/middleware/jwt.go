package middleware

import (
	"errors"
	"net/http"
	"strings"

	"backend/pkg/config"
	customeErr "backend/pkg/errors"
	"backend/pkg/jwt"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	secret string
}

func CreateJWTMiddleware(cfg *config.Config) *JWTMiddleware {
	return &JWTMiddleware{
		secret: cfg.Security.SecretKey,
	}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过登录、登出、创建用户等公开接口
		path := c.Request.URL.Path
		if path == "/api/v1/user/login" ||
			path == "/api/v1/user/create" ||
			path == "/api/v1/ping" {
			c.Next()
			return
		}

		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, response.Response{
				Code:    customeErr.DefaultError,
				Data:    nil,
				Message: "missing authorization header",
			})
			c.Abort()
			return
		}

		// 验证 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, response.Response{
				Code:    customeErr.DefaultError,
				Data:    nil,
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 JWT
		claims, err := jwt.ValidateToken(tokenString, m.secret)
		if err != nil {
			message := "invalid token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				message = "token has expired"
			}
			response.Error(c, http.StatusUnauthorized, response.Response{
				Code:    customeErr.DefaultError,
				Data:    nil,
				Message: message,
			})
			c.Abort()
			return
		}

		// 将用户信息放入 Context
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (userId int64, username string, ok bool) {
	userIdVal, exist1 := c.Get("userId")
	usernameVal, exist2 := c.Get("username")
	if !exist1 || !exist2 {
		return 0, "", false
	}
	userId, ok = userIdVal.(int64)
	if !ok {
		return 0, "", false
	}
	username, ok = usernameVal.(string)
	return
}
