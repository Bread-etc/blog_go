package middleware

import (
	"net/http"
	"strings"

	"go-blog/pkg/errs"
	jwtpkg "go-blog/pkg/jwt"
	"go-blog/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			abortUnauthorized(c, "authorization header is required")
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abortUnauthorized(c, "invalid token format")
			return
		}

		claims, err := jwtpkg.ParseToken(parts[1])
		if err != nil {
			abortUnauthorized(c, "invalid or expired token")
			return
		}
		if claims.UserID == "" || claims.Username == "" || claims.Role == "" {
			abortUnauthorized(c, "invalid token claims")
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	response.Error(c, errs.New(http.StatusUnauthorized, errs.CodeUnauthorized, message))
	c.Abort()
}
