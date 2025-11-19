package controller

import (
	"net/http"

	jwtpkg "go-blog/pkg/jwt"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.IUserService
}

// NewUserController 接口一个接口类型
func NewUserController(userService service.IUserService) *UserController {
	return &UserController{UserService: userService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 方法处理用户登录
func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := uc.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to genreate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// GetProfile 方法获取用户信息
func (uc *UserController) GetProfile(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("userID")
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
