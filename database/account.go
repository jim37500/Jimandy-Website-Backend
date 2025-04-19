package database

import "Jimandy-Website-Backend/model"

// 依 主鍵 取得 帳號
func GetAccountByID(id uint64) (account model.Account) {
	db.First(&account, id)

	return
}

// 依 電子郵件 取得 帳號
func GetAccountByEmail(email string) (account model.Account) {
	db.Where("email = ?", email).First(&account)

	return 
}

// 新增帳號
func AddAccount(account model.Account) bool {
	return db.Save(&account).Error == nil
}
