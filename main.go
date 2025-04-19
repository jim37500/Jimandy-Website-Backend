package main

import (
	"log"
	"os"
	"path/filepath"
	"Jimandy-Website-Backend/configuration"
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/router"
)

func main() {
	// 設定執行檔路徑
	ExecutPath, _ := os.Executable()
	ExecutPath = filepath.Dir(ExecutPath)
	configuration.ExecutPath = ExecutPath

	configuration.ReadConfiguration() // 讀取設定檔
	
	log.Println("Opening Project DB...")

	// 連線資料庫
	database.Open()

	log.Println("Starting Project...")

	// 監聽網頁服務
	router.Run()

	log.Println("Stoping Project...")
}