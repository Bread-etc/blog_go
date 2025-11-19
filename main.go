package main

import (
	"fmt"
	"log"

	"go-blog/config"
	"go-blog/model"
	"go-blog/pkg/database"
	jwtpkg "go-blog/pkg/jwt"
	router "go-blog/router"
	service "go-blog/services"
)

func main() {
	// 加载配置
	config.InitConfig()
	// 初始化数据库连接
	db, err := database.InitMySQL()
	if err != nil {
		log.Fatalf("❌ Failed to connect the database: %v", err)
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
		log.Fatalf("❌ Data table migration failed: %v", err)
	}
	log.Println("✅ Data table migration successfully!")

	// 初始化 JWT
	jcfg := &jwtpkg.Config{
		Algorithm:      config.AppConfig.JWT.Algorithm,
		Secret:         config.AppConfig.JWT.Secret,
		PrivateKeyPath: config.AppConfig.JWT.PrivateKeyPath,
		PublicKeyPath:  config.AppConfig.JWT.PublicKeyPath,
		ExpireHours:    config.AppConfig.JWT.ExpireHours,
	}
	if err := jwtpkg.Init(jcfg); err != nil {
		log.Fatalf("❌ Failed to init JWT: %v", err)
	}
	log.Println("✅ JWT initialized")

	// 初始化 Service 并检查 / 创建默认管理员
	userService := service.NewUserService(db)
	if err := userService.CreateAdminIfNotExists(); err != nil {
		log.Printf("❌ Failed to create default adminadministrator: %v", err)
	} else {
		log.Println("✅ Default adminadministrator created successfully!")
	}
	r := router.InitRouter(db)
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", port)
	log.Printf("🚀 Server started at: http://localhost%s successfully!", addr)
	r.Run(addr)
}
