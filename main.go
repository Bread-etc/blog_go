package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	logger.InitLogger("logs/server.log", "info")
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

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 Goroutine 中启动服务器
	go func() {
		logger.Log.Infof("🚀 Server started at: http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("❌ Listen: %s\n", err)
		}
	}()

	// 优雅关闭（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Infof("⛔️ Shutting down server...")

	// 创建一个 5 秒超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown 会等待活跃连接完成，然后关闭
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("❌ Server Shutdown (Force): %s", err)
	}

	logger.Log.Infof("✅ Server exiting")
}
