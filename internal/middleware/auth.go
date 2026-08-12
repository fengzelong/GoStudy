package middleware

import (
	"net/http"
	"strings"

	"GoStudy/internal/auth"
	"GoStudy/internal/domain"
	"GoStudy/internal/response"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"
const UserRoleKey = "user_role"

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
		c.Set(UserRoleKey, claims.Role)
		c.Next()
	}
}

// RequireRole 限制只有指定角色可以访问路由。
func RequireRole(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get(UserRoleKey)
		currentRole, _ := value.(domain.UserRole)
		if currentRole != role {
			response.Error(c, http.StatusForbidden, "insufficient permission")
			c.Abort()
			return
		}
		c.Next()
	}
}
