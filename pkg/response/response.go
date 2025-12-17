package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构体
type Response struct {
	Code    int    `json:"code"`    // 业务状态码
	Message string `json:"message"` // 提示信息
	Data    any    `json:"data"`    // 数据
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: msg,
		Data:    nil,
	})
}

// ErrorWithCode 自定义业务码的错误响应
func ErrorWithCode(c *gin.Context, httpCode int, businessCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    businessCode,
		Message: msg,
		Data:    nil,
	})
}
