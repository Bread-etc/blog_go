package database

import (
	"fmt"
	"go-blog/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化 mysql 数据库
func InitMySQL() (*gorm.DB, error) {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Charset,
		dbConfig.ParseTime,
		dbConfig.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Database connection successfully!")
	return db, nil
}
