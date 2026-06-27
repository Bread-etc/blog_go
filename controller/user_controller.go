package controller

import (
	"net/http"

	"go-blog/dto"
	"go-blog/pkg/crypto"
	"go-blog/pkg/errs"
	jwtpkg "go-blog/pkg/jwt"
	"go-blog/pkg/response"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{UserService: userService}
}

// GetPublicKey 获取公钥
func (uc *UserController) GetPublicKey(c *gin.Context) {
	publicKey, err := uc.UserService.GetPublicKey()
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.PublicKeyResp{
		PublicKey: publicKey,
	})
}

// Login 登录
func (uc *UserController) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	// 解密传入的 RSA 加密后的 Base64 字符串
	plainPassword, err := crypto.Decrypt(req.Password)
	if err != nil {
		response.Error(c, errs.New(http.StatusBadRequest, errs.CodeInvalidPasswordEncrypt, "invalid password encryption"))
		return
	}

	// 传入解密后的 plainPassword
	user, err := uc.UserService.AuthenticateUser(req.Username, plainPassword)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 利用 jwt 发放 token
	token, err := jwtpkg.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Error(c, errs.Wrap(http.StatusInternalServerError, errs.CodeInternalError, "failed to generate token", err))
		return
	}

	response.Success(c, dto.LoginResp{
		Token: token,
		User:  *user,
	})
}

// GetProfile 获取当前登录用户信息
func (uc *UserController) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	username := c.GetString("username")
	role := c.GetString("role")

	if userID == "" || username == "" || role == "" {
		response.Error(c, errs.New(http.StatusUnauthorized, errs.CodeUnauthorized, "invalid token claims"))
		return
	}

	response.Success(c, dto.UserBrief{
		ID:       userID,
		Username: username,
		Role:     role,
	})
}

// ChangePassword 修改密码
func (uc *UserController) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, errs.New(http.StatusUnauthorized, errs.CodeUnauthorized, "invalid token claims"))
		return
	}

	oldPlainPassword, err := crypto.Decrypt(req.OldPassword)
	if err != nil {
		response.Error(c, errs.New(http.StatusBadRequest, errs.CodeInvalidPasswordEncrypt, "invalid password encryption"))
		return
	}

	newPlainPassword, err := crypto.Decrypt(req.NewPassword)
	if err != nil {
		response.Error(c, errs.New(http.StatusBadRequest, errs.CodeInvalidPasswordEncrypt, "invalid password encryption"))
		return
	}

	if err := uc.UserService.ChangePassword(userID, oldPlainPassword, newPlainPassword); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}
