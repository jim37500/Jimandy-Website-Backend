package model

import (
	"gorm.io/gorm"
)

// 自動遷移資料庫
func AutoMigrate(db *gorm.DB) {
	migrateTable(db, &Account{})
	migrateTable(db, &Athlete{})
	migrateTable(db, &Activity{})
	migrateTable(db, &Lap{})

	checkTableData(db)
}

// 轉移資料表結構
func migrateTable(db *gorm.DB, structure interface{}) {
	if !db.Migrator().HasTable(structure) {
		_ = db.Migrator().CreateTable(structure)
	}
	_ = db.AutoMigrate(structure)
}

// 檢查資料表有無資料
func checkTableData(db *gorm.DB) {
	// 若沒帳號
	var account Account
	if db.First(&account).RowsAffected == 0 {
		myAccount := Account{Name: "管理員", Email: "admin@gmail.com", Password: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", Status: 1, IsAdmin: 1}
		db.Create(&myAccount)
	}
}
