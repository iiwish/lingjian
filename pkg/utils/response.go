package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
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
