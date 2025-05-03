package model

import "gorm.io/gorm"

const (
	Disabled int8 = iota
	Enabled
	UnActive
)

type Account struct {
	gorm.Model

	Name         string `gorm:"comment:名稱"`
	Email        string `gorm:"unique;not null;size:64;comment:電子郵件"`
	Status       int8   `gorm:"index;comment:狀態 0停用 1啟用;default:1"`
	IsAdmin      int    `gorm:"comment:是否為管理員 1是 0否"`
}
