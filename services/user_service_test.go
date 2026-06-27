package service

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"go-blog/model"
	"go-blog/pkg/crypto"
	"go-blog/pkg/errs"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))

	return db
}

func createTestUser(t *testing.T, db *gorm.DB, username string, password string, role string) model.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := model.User{
		Username: username,
		Password: string(hash),
		Role:     role,
	}
	require.NoError(t, db.Create(&user).Error)

	return user
}

func assertUserAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func TestUserService_CreateAdminIfNotExists(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	require.NoError(t, svc.CreateAdminIfNotExists())

	var admin model.User
	require.NoError(t, db.Where("role = ?", defaultAdminRole).First(&admin).Error)
	assert.NotEmpty(t, admin.ID)
	assert.Equal(t, defaultAdminUsername, admin.Username)
	assert.Equal(t, defaultAdminRole, admin.Role)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(defaultAdminPassword)))

	require.NoError(t, svc.CreateAdminIfNotExists())

	var count int64
	require.NoError(t, db.Model(&model.User{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestUserService_CreateAdminIfNotExists_DatabaseError(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.CreateAdminIfNotExists()
	require.Error(t, err)
	assertUserAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestUserService_AuthenticateUser(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	user := createTestUser(t, db, "testuser", "secret123", "admin")

	got, err := svc.AuthenticateUser("testuser", "secret123")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, "testuser", got.Username)
	assert.Equal(t, "admin", got.Role)
}

func TestUserService_AuthenticateUser_InvalidPassword(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)
	createTestUser(t, db, "testuser", "secret123", "admin")

	user, err := svc.AuthenticateUser("testuser", "wrongpassword")
	require.Error(t, err)
	assert.Nil(t, user)
	assertUserAppError(t, err, http.StatusUnauthorized, errs.CodeInvalidCredentials)
}

func TestUserService_AuthenticateUser_UserNotFound(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	user, err := svc.AuthenticateUser("ghost", "secret123")
	require.Error(t, err)
	assert.Nil(t, user)
	assertUserAppError(t, err, http.StatusUnauthorized, errs.CodeInvalidCredentials)
}

func TestUserService_AuthenticateUser_DatabaseError(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	user, err := svc.AuthenticateUser("testuser", "secret123")
	require.Error(t, err)
	assert.Nil(t, user)
	assertUserAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestUserService_ChangePassword(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)
	user := createTestUser(t, db, "changer", "oldpass", "admin")

	require.NoError(t, svc.ChangePassword(user.ID, "oldpass", "newpass_secure"))

	var updatedUser model.User
	require.NoError(t, db.First(&updatedUser, "id = ?", user.ID).Error)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpass_secure")))
	assert.Error(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("oldpass")))
}

func TestUserService_ChangePassword_IncorrectOldPassword(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)
	user := createTestUser(t, db, "changer", "oldpass", "admin")

	err := svc.ChangePassword(user.ID, "wrongold", "newpass_secure")
	require.Error(t, err)
	assertUserAppError(t, err, http.StatusBadRequest, errs.CodeIncorrectOldPassword)
}

func TestUserService_ChangePassword_UserNotFound(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	err := svc.ChangePassword("missing-user-id", "oldpass", "newpass_secure")
	require.Error(t, err)
	assertUserAppError(t, err, http.StatusNotFound, errs.CodeNotFound)
}

func TestUserService_ChangePassword_DatabaseError(t *testing.T) {
	db := setupUserTestDB(t)
	svc := NewUserService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.ChangePassword("user-id", "oldpass", "newpass_secure")
	require.Error(t, err)
	assertUserAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestUserService_GetPublicKey(t *testing.T) {
	require.NoError(t, crypto.InitRSAKeyPair())

	svc := NewUserService(nil)
	publicKey, err := svc.GetPublicKey()
	require.NoError(t, err)
	assert.True(t, strings.Contains(publicKey, "BEGIN PUBLIC KEY"))
	assert.True(t, strings.Contains(publicKey, "END PUBLIC KEY"))
}
