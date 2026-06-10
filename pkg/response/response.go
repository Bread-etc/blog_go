package response

import (
	"go-blog/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 基础响应结构体
type Response struct {
	Code      int    `json:"code"`      // HTTP 状态码
	ErrorCode string `json:"errorCode"` // 业务错误码
	Message   string `json:"message"`   // 提示信息
	Data      any    `json:"data"`      // 数据
}

// 成功响应封装
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		ErrorCode: errs.CodeSuccess,
		Message:   "success",
		Data:      data,
	})
}

// 错误响应封装
func Error(c *gin.Context, appErr *errs.Error) {
	if appErr == nil {
		appErr = errs.New(http.StatusInternalServerError, errs.CodeInternalError, "internal server error")
	}

	c.JSON(appErr.HTTPStatus, Response{
		Code:      appErr.HTTPStatus,
		ErrorCode: appErr.Code,
		Message:   appErr.Message,
		Data:      nil,
	})
}

// 自定义业务错误码的错误响应封装
func ErrorFrom(c *gin.Context, err error) {
	httpStatus, errorCode, message := errs.From(err)

	c.JSON(httpStatus, Response{
		Code:      httpStatus,
		ErrorCode: errorCode,
		Message:   message,
		Data:      nil,
	})
}
