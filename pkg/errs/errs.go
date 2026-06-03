package errs

import (
	"errors"
	"net/http"
)

type Error struct {
	HTTPStatus int    // HTTP 状态码
	Code       string // 业务错误码
	Message    string // 响应消息
	Cause      error  // 底层原始错误
}

// 创建业务错误
func New(httpStatus int, code string, message string) *Error {
	return &Error{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// 创建业务错误，保留底层错误原因
func Wrap(httpStatus int, code string, message string, cause error) *Error {
	return &Error{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Cause:      cause,
	}
}

// 为 *Error 类型定义 Error 方法，返回 string
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// 为 *Error 类型定义 Unwrap 方法，返回 error
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func From(err error) (int, string, string) {
	if err == nil {
		return http.StatusOK, CodeSuccess, "success"
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus, appErr.Code, appErr.Message
	}

	return http.StatusInternalServerError, CodeInternalError, "internal server error"
}

func DefaultCode(httpStatus int) string {
	switch httpStatus {
	case http.StatusBadRequest:
		return CodeInvalidParams
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	default:
		return CodeInternalError
	}
}
