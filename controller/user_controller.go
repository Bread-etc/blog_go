package controller

import (
	"net/http"

	jwtpkg "go-blog/pkg/jwt"
	"go-blog/pkg/response"
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
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := uc.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to genreate token")
		return
	}

	response.Success(c, gin.H{
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
	username := c.GetString("username") // 暂时未使用的变量，先注释掉或删除
	userID := c.GetString("userID")
	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码接口
func (uc *UserController) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")

	if err := uc.UserService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Password updated successfully"})
}
