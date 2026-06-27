package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/crypto"
	"go-blog/pkg/errs"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin"
	defaultAdminRole     = "admin"
)

type IUserService interface {
	AuthenticateUser(username, password string) (*dto.UserBrief, error)
	CreateAdminIfNotExists() error
	ChangePassword(userID, oldPassword, newPassword string) error
	GetPublicKey() (string, error)
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

var _ IUserService = (*UserService)(nil)

// AuthenticateUser 验证用户名 + 密码
func (us *UserService) AuthenticateUser(username, password string) (*dto.UserBrief, error) {
	var user model.User

	if err := us.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusUnauthorized, errs.CodeInvalidCredentials, "invalid username or password")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to authenticate user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errs.New(http.StatusUnauthorized, errs.CodeInvalidCredentials, "invalid username or password")
	}

	return toUserBrief(user), nil
}

// CreateAdminIfNotExists 用于初始化默认管理员
func (us *UserService) CreateAdminIfNotExists() error {
	var count int64

	if err := us.DB.Model(&model.User{}).Where("role = ?", defaultAdminRole).Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count admin users", err)
	}

	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeInternalError, "failed to hash default admin password", err)
	}

	admin := model.User{
		Username: defaultAdminUsername,
		Password: string(hash),
		Role:     defaultAdminRole,
	}

	if err := us.DB.Create(&admin).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create default admin", err)
	}

	return nil
}

// ChangePassword 修改密码
func (us *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	var user model.User

	// 查找用户
	if err := us.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeNotFound, "user not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get user", err)
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errs.New(http.StatusBadRequest, errs.CodeIncorrectOldPassword, "incorrect old password")
	}

	// 加密新密码
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeInternalError, "failed to hash new password", err)
	}

	// 更新数据库
	if err := us.DB.Model(&user).Update("password", string(newHash)).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update password", err)
	}

	return nil
}

// GetPublicKey 获取加密公钥
func (us *UserService) GetPublicKey() (string, error) {
	publicKey, err := crypto.GetPublicKey()

	if err != nil {
		return "", errs.Wrap(http.StatusInternalServerError, errs.CodeInternalError, "failed to get public key", err)
	}

	return publicKey, nil
}

func toUserBrief(user model.User) *dto.UserBrief {
	return &dto.UserBrief{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}
