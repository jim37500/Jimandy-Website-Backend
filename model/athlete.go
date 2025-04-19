package model

// 運動員
type Athlete struct {
	ID                uint64   `gorm:"primarykey"`
	AccountID         uint     `gorm:"comment:帳號主鍵"`
	ClientID          uint     `gorm:"comment:用戶端 ID"`
	ClientSecret      string   `gorm:"comment:用戶端密碼"`
	AuthorizationCode string   `gorm:"comment:授權代碼"`
	AccessToken       string   `gorm:"comment:通行權杖"`
	RefreshToken      string   `gorm:"comment:換發權杖"`
	RunLap            Lap      `gorm:"foreignKey:AthleteID"`
	RunActivity       Activity `gorm:"foreignKey:AthleteID"`
}
