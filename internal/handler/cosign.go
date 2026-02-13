package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sm2-cosign/backend/internal/middleware"
	"github.com/sm2-cosign/backend/internal/service"
	"github.com/sm2-cosign/backend/pkg/response"
)

// CosignHandler 协同签名处理器
type CosignHandler struct {
	cosignService *service.CosignService
}

// NewCosignHandler 创建协同签名处理器实例
func NewCosignHandler() *CosignHandler {
	return &CosignHandler{
		cosignService: service.NewCosignService(),
	}
}

// KeyInit 密钥初始化
// @Summary 密钥初始化
// @Description 初始化或重新生成SM2协同密钥对
// @Tags 协同签名
// @Accept json
// @Produce json
// @Param request body service.KeyInitRequest true "密钥初始化请求"
// @Success 200 {object} response.Response{data=service.KeyInitResponse}
// @Router /api/key/init [post]
func (h *CosignHandler) KeyInit(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return response.Error(c, response.CodeUnauthorized)
	}

	var req service.KeyInitRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	// 使用当前登录用户的ID
	req.UserID = userID

	result, code := h.cosignService.KeyInit(&req, c.IP())
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, result)
}

// Sign 协同签名
// @Summary 协同签名
// @Description 执行SM2协同签名
// @Tags 协同签名
// @Accept json
// @Produce json
// @Param request body service.SignRequest true "签名请求"
// @Success 200 {object} response.Response{data=service.SignResponse}
// @Router /api/sign [post]
func (h *CosignHandler) Sign(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return response.Error(c, response.CodeUnauthorized)
	}

	var req service.SignRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	// 使用当前登录用户的ID
	req.UserID = userID

	result, code := h.cosignService.Sign(&req, c.IP())
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, result)
}

// Decrypt 协同解密
// @Summary 协同解密
// @Description 执行SM2协同解密
// @Tags 协同签名
// @Accept json
// @Produce json
// @Param request body service.DecryptRequest true "解密请求"
// @Success 200 {object} response.Response{data=service.DecryptResponse}
// @Router /api/decrypt [post]
func (h *CosignHandler) Decrypt(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return response.Error(c, response.CodeUnauthorized)
	}

	var req service.DecryptRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	// 使用当前登录用户的ID
	req.UserID = userID

	result, code := h.cosignService.Decrypt(&req, c.IP())
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, result)
}
