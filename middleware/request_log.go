package middleware

import (
	"net/http"
	"time"

	"go-blog/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLog 接收 Gin 默认的日志并用 zap 记录
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// 生成 uuid 并设置响应头
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		if logger.Log == nil {
			return
		}

		status := c.Writer.Status()
		latency := time.Since(startTime)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		fields := []any{
			"requestId", requestID,
			"method", c.Request.Method,
			"status", status,
			"latency", latency.String(),
			"path", path,
			"uri", c.Request.RequestURI,
			"clientIP", c.ClientIP(),
			"userAgent", c.Request.UserAgent(),
			"userID", c.GetString("userID"),
		}
		// 记录业务代码主动抛出 Gin context 的错误
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Log.Errorw("http request", fields...)
		case status >= http.StatusBadRequest:
			logger.Log.Warnw("http request", fields...)
		default:
			logger.Log.Infow("http request", fields...)
		}
	}
}
