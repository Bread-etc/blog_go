package main

import (
	"fmt"

	"go-blog/config"
	"go-blog/model"
	crypto "go-blog/pkg/crypto"
	"go-blog/pkg/database"
	jwtpkg "go-blog/pkg/jwt"
	"go-blog/pkg/logger"
	router "go-blog/router"
	service "go-blog/services"
)

func main() {
	// 加载配置
	config.InitConfig()
	logger.InitLogger("logs/server.log", "debug")
	// 初始化数据库连接
	db, err := database.InitMySQL()
	if err != nil {
		logger.Log.Errorf("❌ Failed to connect the database: %v", err)
	}

	// 自动迁移模型
	err = db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Post{},
		&model.SiteConfig{},
	)
	if err != nil {
		logger.Log.Errorf("❌ Data table migration failed: %v", err)
	}
	logger.Log.Infof("✅ Data table migration successfully!")

	// 初始化 JWT
	jcfg := &jwtpkg.Config{
		Algorithm:      config.AppConfig.JWT.Algorithm,
		Secret:         config.AppConfig.JWT.Secret,
		PrivateKeyPath: config.AppConfig.JWT.PrivateKeyPath,
		PublicKeyPath:  config.AppConfig.JWT.PublicKeyPath,
		ExpireHours:    config.AppConfig.JWT.ExpireHours,
	}
	if err := jwtpkg.Init(jcfg); err != nil {
		logger.Log.Errorf("❌ Failed to init JWT: %v", err)
	}
	logger.Log.Infof("✅ JWT initialized successfully!")

	// 初始化 RSA 密钥对
	if err := crypto.InitRSAKeyPair(); err != nil {
		logger.Log.Errorf("❌ Failed to init RSA KeyPair: %v", err)
	}
	logger.Log.Infof("✅ RSA KeyPair initialized sucessfully!")

	// 初始化 Service 并检查 / 创建默认管理员
	userService := service.NewUserService(db)
	if err := userService.CreateAdminIfNotExists(); err != nil {
		logger.Log.Errorf("❌ Failed to create default adminadministrator: %v", err)
	} else {
		logger.Log.Infof("✅ Default administrator created successfully!")
	}

	r := router.InitRouter(db)
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", port)
	logger.Log.Infof("🚀 Server started at: http://localhost%s successfully!", addr)
	r.Run(addr)
}
