package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sm2-cosign/backend/internal/repository"
	"github.com/sm2-cosign/backend/internal/service"
	"github.com/sm2-cosign/backend/pkg/response"
)

// AdminHandler 管理处理器
type AdminHandler struct {
	userService   *service.UserService
	cosignService *service.CosignService
	auditRepo     *repository.AuditLogRepository
}

// NewAdminHandler 创建管理处理器实例
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		userService:   service.NewUserService(),
		cosignService: service.NewCosignService(),
		auditRepo:     repository.NewAuditLogRepository(),
	}
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户列表（分页）
// @Tags 管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /mapi/users [get]
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, code := h.userService.ListUsers(page, pageSize)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, fiber.Map{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} response.Response{data=model.User}
// @Router /mapi/users/{id} [get]
func (h *AdminHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, response.CodeInvalidParam)
	}

	user, code := h.userService.GetUserInfo(id)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} response.Response
// @Router /mapi/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, response.CodeInvalidParam)
	}

	code := h.userService.DeleteUser(id)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, nil)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 启用或禁用用户
// @Tags 管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param request body map[string]int true "状态请求"
// @Success 200 {object} response.Response
// @Router /mapi/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, response.CodeInvalidParam)
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.CodeInvalidParam)
	}

	code := h.userService.UpdateUserStatus(id, req.Status)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, nil)
}

// ListKeys 获取密钥列表
// @Summary 获取密钥列表
// @Description 获取所有密钥列表（分页）
// @Tags 管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /mapi/keys [get]
func (h *AdminHandler) ListKeys(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keys, total, code := h.cosignService.ListKeys(page, pageSize)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, fiber.Map{
		"list":     keys,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteKey 删除密钥
// @Summary 删除密钥
// @Description 删除指定密钥
// @Tags 管理
// @Accept json
// @Produce json
// @Param id path string true "密钥ID"
// @Success 200 {object} response.Response
// @Router /mapi/keys/{id} [delete]
func (h *AdminHandler) DeleteKey(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, response.CodeInvalidParam)
	}

	code := h.cosignService.DeleteKey(id)
	if code != response.CodeSuccess {
		return response.Error(c, code)
	}

	return response.Success(c, nil)
}

// ListLogs 查询审计日志
// @Summary 查询审计日志
// @Description 查询审计日志（分页）
// @Tags 管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param action query string false "操作类型"
// @Param user_id query string false "用户ID"
// @Success 200 {object} response.Response
// @Router /mapi/logs [get]
func (h *AdminHandler) ListLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	action := c.Query("action")
	userID := c.Query("user_id")

	logs, total, err := h.auditRepo.List(page, pageSize, action, userID)
	if err != nil {
		return response.Error(c, response.CodeDBError)
	}

	return response.Success(c, fiber.Map{
		"list":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// Health 健康检查
// @Summary 健康检查
// @Description 服务健康检查
// @Tags 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /mapi/health [get]
func (h *AdminHandler) Health(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"status": "ok",
	})
}

// Stats 系统统计
// @Summary 系统统计
// @Description 获取系统统计数据
// @Tags 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /mapi/stats [get]
func (h *AdminHandler) Stats(c *fiber.Ctx) error {
	// 获取统计数据
	var userCount, keyCount, sessionCount int64
	db := repository.GetDB()
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	db.QueryRow("SELECT COUNT(*) FROM keys").Scan(&keyCount)
	db.QueryRow("SELECT COUNT(*) FROM sessions WHERE expires_at > datetime('now')").Scan(&sessionCount)

	return response.Success(c, fiber.Map{
		"users":    userCount,
		"keys":     keyCount,
		"sessions": sessionCount,
	})
}
