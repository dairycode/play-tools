package database

import (
	"log"
	"play-tools/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("play-tools.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移数据库
	err = DB.AutoMigrate(
		&model.User{},
		&model.GameRoom{},
		&model.RoomMember{},
		&model.Transaction{},
		&model.Settlement{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
}
