package service

import (
	"go-blog/model"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupUserTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	db.AutoMigrate(&model.User{})
	return db
}

func TestUserService_CreateAdminIfNotExists(t *testing.T) {
	db := setupUserTestDB()
	svc := NewUserService(db)

	err := svc.CreateAdminIfNotExists()
	assert.NoError(t, err)

	// 验证管理员是否存在
	var admin model.User
	err = db.Where("role = ?", "admin").First(&admin).Error
	assert.NoError(t, err)
	assert.Equal(t, "admin", admin.Username)

	err = svc.CreateAdminIfNotExists()
	assert.NoError(t, err)

	// 验证重复创建
	var count int64
	db.Model(&model.User{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestUserService_AuthenticateUser(t *testing.T) {
	db := setupUserTestDB()
	svc := NewUserService(db)

	// 准备数据：手动创建一个用户
	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := model.User{
		Username: "testuser",
		Password: string(hashed),
		Role:     "user",
	}
	db.Create(&user)

	// Case 1: 密码正确
	u, err := svc.AuthenticateUser("testuser", password)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", u.Username)

	// Case 2: 密码错误
	_, err = svc.AuthenticateUser("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")

	// Case 3: 用户不存在
	_, err = svc.AuthenticateUser("ghost", password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestUserService_ChangePassword(t *testing.T) {
	db := setupUserTestDB()
	svc := NewUserService(db)

	oldPwd := "oldpass"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(oldPwd), bcrypt.DefaultCost)
	user := model.User{Username: "changer", Password: string(hashed)}
	db.Create(&user)

	// Case 1: 旧密码错误
	err := svc.ChangePassword(user.ID, "wrongold", "newpass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incorrect old password")

	// Case 2: 成功修改
	newPwd := "newpass_secure"
	err = svc.ChangePassword(user.ID, oldPwd, newPwd)
	assert.NoError(t, err)

	// 验证数据库里存的是新密码的哈希
	var updatedUser model.User
	db.First(&updatedUser, "id = ?", user.ID)

	// 验证新密码能通过校验
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPwd))
	assert.NoError(t, err)
}
