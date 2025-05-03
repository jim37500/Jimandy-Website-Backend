package model

import "time"

// AccessToken 存儲用戶的 token 資訊
type AccessToken struct {
	ID          uint      `gorm:"primarykey"`
	AccountID   uint      `gorm:"index;comment:帳號主鍵"`
	AccessToken string    `gorm:"unique;comment:存取權杖"`
	DeviceInfo  string    `gorm:"comment:裝置資訊"`
	CreatedAt   time.Time `gorm:"comment:建立時間"`
	ExpiresAt   time.Time `gorm:"index;comment:過期時間"`
	IsRevoked   bool      `gorm:"default:false;comment:是否已撤銷"`
}

// RefreshToken 存儲用戶的 refresh token 資訊
type RefreshToken struct {
	ID           uint      `gorm:"primarykey"`
	AccountID    uint      `gorm:"index;comment:帳號主鍵"`
	RefreshToken string    `gorm:"unique;comment:刷新權杖"`
	DeviceInfo   string    `gorm:"comment:裝置資訊"`
	CreatedAt    time.Time `gorm:"comment:建立時間"`
	ExpiresAt    time.Time `gorm:"index;comment:過期時間"`
	IsRevoked    bool      `gorm:"default:false;comment:是否已撤銷"`
}
