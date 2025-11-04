package main

import (
	"fmt"
	"go-blog/config"
	"go-blog/model"
	"go-blog/pkg/database"
	"go-blog/router"
	"log"
)

func main() {
	config.InitConfig()
	database.InitMySQL()

	// 自动迁移模型
	err := database.DB.AutoMigrate(
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

	// 初始化路由
	r := router.InitRouter()
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", port)
	log.Printf("🚀 Server started at: http://localhost%s successfully!", addr)
	r.Run(addr)
}
