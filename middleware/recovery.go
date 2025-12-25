package middleware

import (
	"fmt"
	"go-blog/pkg/logger"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 错误恢复中间件 - 捕获 panic 并记录 Stack Trace
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 记录 Panic 原因和堆栈信息
				logger.Log.Errorf("[Recovery] panic recovered:\n%s\n%v\nStack: %s", string(httpRequest), err, string(debug.Stack()))
				// 返回 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
			}
		}()
		c.Next()
	}
}
