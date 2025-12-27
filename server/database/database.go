package database

import (
	"fmt"
	"log"
	"play-tools/config"
	"play-tools/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化数据库
func Init() (err error) {
	// 从配置文件读取数据库配置
	dbConfig := config.AppConfig.Database

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	// 自动迁移数据库
	if err = DB.AutoMigrate(
		&model.User{},
		&model.GameRoom{},
		&model.RoomMember{},
		&model.Transaction{},
		&model.Settlement{},
	); err != nil {
		return
	}

	log.Println("Database initialized successfully")
	return
}
