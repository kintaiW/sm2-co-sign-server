package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sm2-cosign/backend/internal/middleware"
	"github.com/sm2-cosign/backend/internal/service"
	"github.com/sm2-cosign/backend/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户并生成SM2协同密钥对
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册请求"
// @Success 200 {object} response.Response{data=service.RegisterResponse}
// @Router /api/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	result, code := h.userService.Register(&req, c.IP())
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, result)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Router /api/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	result, code := h.userService.Login(&req, c.IP())
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, result)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，删除Token
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// 从 Header 获取 Token
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Success(c, nil)
	}

	parts := []byte(authHeader)
	if len(parts) > 7 && string(parts[:7]) == "Bearer " {
		token := string(parts[7:])
		h.userService.Logout(token)
	}

	return response.Success(c, nil)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/user/info [get]
func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return response.Error(c, response.CodeUnauthorized)
	}

	user, code := h.userService.GetUserInfo(userID)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, user)
}
