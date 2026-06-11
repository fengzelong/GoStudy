package middleware

import (
	"net/http"
	"strings"

	"GoStudy/internal/auth"
	"GoStudy/internal/response"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// Auth 校验 Bearer Token，并把用户 ID 放入 Gin 上下文供后续处理使用。
func Auth(tokens *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		claims, err := tokens.Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.Subject)
		c.Next()
	}
}
