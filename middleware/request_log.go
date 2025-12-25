package middleware

import (
	"go-blog/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLog 接收 Gin 默认的日志并用 zap 记录
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		duration := time.Since(startTime)
		stauts := c.Writer.Status()
		// 记录日志
		logger.Log.Infof("[HTTP] %d | %12v | %15s | %-7s %s",
			stauts,
			duration,
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI)
	}
}
