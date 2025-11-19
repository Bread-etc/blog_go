package service

import (
	"errors"
	"go-blog/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// IUserService 定义用户服务接口
type IUserService interface {
	AuthenticateUser(username, password string) (*model.User, error)
	CreateAdminIfNotExists() error
	ChangePassword(userID, oldPassword, newPassword string) error
}

type UserService struct {
	DB *gorm.DB
}

// NewUserService 实例化 UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// 确保 UserService 实现了 IUserService 接口
var _ IUserService = (*UserService)(nil)

// AuthenticateUser 验证用户名/密码
func (us *UserService) AuthenticateUser(username, password string) (*model.User, error) {
	var user model.User

	if err := us.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

// CreateAdminIfNotExists 用于初始化默认管理员
func (us *UserService) CreateAdminIfNotExists() error {
	var count int64
	if err := us.DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		pw := "admin"
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin := model.User{
			ID:       uuid.NewString(),
			Username: "admin",
			Password: string(hash),
			Role:     "admin",
		}
		return us.DB.Create(&admin).Error
	}
	return nil
}

// ChangePassword 修改密码
func (us *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	var user model.User
	// 查找用户
	if err := us.DB.First(&user, "id = ?", userID).Error; err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// 加密新密码
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新数据库
	return us.DB.Model(&user).Update("password", string(newHash)).Error
}
