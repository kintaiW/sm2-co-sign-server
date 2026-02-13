package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sm2-cosign/backend/internal/service"
	"github.com/sm2-cosign/backend/pkg/response"
)

const (
	// ContextKeyUserID 用户ID上下文键
	ContextKeyUserID = "user_id"
	// ContextKeySession 会话上下文键
	ContextKeySession = "session"
)

// AuthMiddleware Token认证中间件
func AuthMiddleware() fiber.Handler {
	userService := service.NewUserService()

	return func(c *fiber.Ctx) error {
		// 从 Header 获取 Token
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, response.CodeUnauthorized)
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, response.CodeUnauthorized)
		}
		token := parts[1]

		// 验证会话
		session, code := userService.ValidateSession(token)
		if code != response.CodeSuccess {
			return response.Error(c, code)
		}

		// 将用户信息存入上下文
		c.Locals(ContextKeyUserID, session.UserID)
		c.Locals(ContextKeySession, session)

		return c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetSession 从上下文获取会话
func GetSession(c *fiber.Ctx) interface{} {
	return c.Locals(ContextKeySession)
}
