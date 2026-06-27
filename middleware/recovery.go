package middleware

import (
	"net/http"
	"runtime/debug"

	"go-blog/pkg/errs"
	"go-blog/pkg/logger"
	"go-blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery 错误恢复中间件 - 捕获 panic 并记录 Stack Trace
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 Panic 原因和堆栈信息
				if logger.Log != nil {
					path := c.FullPath()
					if path == "" {
						path = c.Request.URL.Path
					}

					logger.Log.Errorw("panic recovered",
						"requestId", c.Writer.Header().Get("X-Request-ID"),
						"method", c.Request.Method,
						"path", path,
						"uri", c.Request.RequestURI,
						"clientIP", c.ClientIP(),
						"userID", c.GetString("userID"),
						"panic", err,
						"stack", string(debug.Stack()),
					)
				}
				// 如果响应没有写出去，才返回统一 500 JSON
				if !c.Writer.Written() {
					response.Error(c, errs.New(http.StatusInternalServerError, errs.CodeInternalError, "internal server error"))
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
