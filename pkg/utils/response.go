package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// CustomError 自定义错误
type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

// NewError 创建新的自定义错误
func NewError(message string) error {
	return &CustomError{message: message}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PaginationResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	// 根据业务错误码设置对应的HTTP状态码
	var httpStatus int
	switch code {
	case 400:
		httpStatus = http.StatusBadRequest
	case 401:
		httpStatus = http.StatusUnauthorized
	case 403:
		httpStatus = http.StatusForbidden
	case 404:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	var httpStatus int
	switch code {
	case 400:
		httpStatus = http.StatusBadRequest
	case 401:
		httpStatus = http.StatusUnauthorized
	case 403:
		httpStatus = http.StatusForbidden
	case 404:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ValidationError 参数验证错误响应
func ValidationError(c *gin.Context, err error) {
	Error(c, 400, fmt.Sprintf("参数验证错误: %v", err))
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, err error) {
	Error(c, 500, fmt.Sprintf("服务器错误: %v", err))
}

// NotFoundError 资源不存在错误响应
func NotFoundError(c *gin.Context, resource string) {
	Error(c, 404, fmt.Sprintf("%s不存在", resource))
}

// UnauthorizedError 未授权错误响应
func UnauthorizedError(c *gin.Context) {
	Error(c, 401, "未授权访问")
}

// ForbiddenError 禁止访问错误响应
func ForbiddenError(c *gin.Context) {
	Error(c, 403, "禁止访问")
}
