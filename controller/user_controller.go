package controller

import (
	"fmt"
	"net/http"

	"go-blog/dto"
	"go-blog/pkg/crypto"
	jwtpkg "go-blog/pkg/jwt"
	"go-blog/pkg/logger"
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

// GetPublicKey 获取公钥接口
func (uc *UserController) GetPublicKey(c *gin.Context) {
	pubKey, err := uc.UserService.GetPublicKey()
	if err != nil {
		logger.Log.Errorf("GetPublicKey service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get public key: %v", err))
		return
	}

	response.Success(c, gin.H{"public_key": pubKey})
}

// Login 登录接口
func (uc *UserController) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("Login bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid input: %v", err))
		return
	}

	// 解密传入的 RSA 加密后的 Base64 字符串
	plainPassword, err := crypto.Decrypt(req.Password)
	if err != nil {
		logger.Log.Warnf("Decrypt password failed: %v", err)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid password encryption: %v", err))
		return
	}

	// 传入解密后的 plainPassword
	user, err := uc.UserService.AuthenticateUser(req.Username, plainPassword)
	if err != nil {
		logger.Log.Errorf("Login service failed: %v", err)
		response.Error(c, http.StatusUnauthorized, fmt.Sprintf("Invalid username or password: %v", err))
		return
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Log.Errorf("Generate token failed: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to genreate token: %v", err))
		return
	}

	response.Success(c, dto.LoginResp{
		Token: token,
		User: dto.UserBrief{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	})
}

// GetProfile 方法获取用户信息
func (uc *UserController) GetProfile(c *gin.Context) {
	username := c.GetString("username")
	userID := c.GetString("userID")

	response.Success(c, dto.UserBrief{
		ID:       userID,
		Username: username,
	})
}

// ChangePassword 修改密码接口
func (uc *UserController) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("ChangePassword bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")

	oldPlain, err := crypto.Decrypt(req.OldPassword)
	if err != nil {
		logger.Log.Warnf("Decrypt password failed: %v", err)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid password encryption: %v", err))
		return
	}

	newPlain, err := crypto.Decrypt(req.NewPassword)
	if err != nil {
		logger.Log.Warnf("Decrypt password failed: %v", err)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid password encryption: %v", err))
		return
	}

	if err := uc.UserService.ChangePassword(userID, oldPlain, newPlain); err != nil {
		logger.Log.Errorf("ChangePassword service error: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Password updated successfully"})
}
