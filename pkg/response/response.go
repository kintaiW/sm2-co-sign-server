package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Code 错误码类型
type Code int

// 错误码定义
const (
	CodeSuccess         Code = 0
	CodeInvalidParam    Code = 10001
	CodeUserExists      Code = 10002
	CodeUserNotFound    Code = 10003
	CodePasswordError   Code = 10004
	CodeTokenInvalid    Code = 10005
	CodeTokenExpired    Code = 10006
	CodeUserDisabled    Code = 10007
	CodeKeyNotFound     Code = 10008
	CodeCryptoError     Code = 10009
	CodeDBError         Code = 10010
	CodeInternalError   Code = 10011
	CodeUnauthorized    Code = 10012
	CodeForbidden       Code = 10013
)

// 错误码消息映射
var codeMessages = map[Code]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "参数错误",
	CodeUserExists:      "用户名已存在",
	CodeUserNotFound:    "用户不存在",
	CodePasswordError:   "密码错误",
	CodeTokenInvalid:    "Token无效",
	CodeTokenExpired:    "Token已过期",
	CodeUserDisabled:    "用户已禁用",
	CodeKeyNotFound:     "密钥不存在",
	CodeCryptoError:     "密码计算错误",
	CodeDBError:         "数据库错误",
	CodeInternalError:   "内部错误",
	CodeUnauthorized:    "未授权",
	CodeForbidden:       "禁止访问",
}

// Response 统一响应结构
type Response struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Code:    CodeSuccess,
		Message: codeMessages[CodeSuccess],
		Data:    data,
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code Code) error {
	message, ok := codeMessages[code]
	if !ok {
		message = codeMessages[CodeInternalError]
	}
	return c.Status(http.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *fiber.Ctx, code Code, message string) error {
	return c.Status(http.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *fiber.Ctx, code Code, data interface{}) error {
	message, ok := codeMessages[code]
	if !ok {
		message = codeMessages[CodeInternalError]
	}
	return c.Status(http.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// GetMessage 获取错误码对应的消息
func GetMessage(code Code) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return codeMessages[CodeInternalError]
}
