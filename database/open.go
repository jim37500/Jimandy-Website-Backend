package database

import (
	"fmt"
	"Jimandy-Website-Backend/configuration"
	"Jimandy-Website-Backend/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Open() {
	var err error
	db, err = gorm.Open(postgres.Open(configuration.Connectionstring), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		time.Sleep(time.Second)
		Open()
		return
	}

	fmt.Println("Connect to Postgres!")

	sqlDB, _ := db.DB()

	// 連線池中空閒連線的最大數量
	sqlDB.SetMaxIdleConns(10)

	// 資料庫連線的最大數量。
	sqlDB.SetMaxOpenConns(100)

	// 連線最長可複用的時間
	sqlDB.SetConnMaxLifetime(time.Hour)

	model.AutoMigrate(db)
}
