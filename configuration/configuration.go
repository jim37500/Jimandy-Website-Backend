package configuration

import (
	"github.com/spf13/viper"
)

var (
	Connectionstring string // 資料庫連線字串
	ExecutPath       string // 執行檔路徑
	JWTKey           []byte // 權杖金鑰
)

// 讀取設定檔
func ReadConfiguration() {
	viper.AutomaticEnv()

	viper.AddConfigPath(ExecutPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig()

	Connectionstring = viper.GetString("CONNECTIONSTRING")
	JWTKey = []byte(viper.GetString("KEY"))
}
