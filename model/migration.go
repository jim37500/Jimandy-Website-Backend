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
	migrateTable(db, &AccessToken{})
	migrateTable(db, &RefreshToken{})

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
}
