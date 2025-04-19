package api

import (
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/model"

	"github.com/gofiber/fiber/v2"
)

// 新增帳號
func AddAccount(context *fiber.Ctx) error {
	var myAccount model.Account
	_ = context.BodyParser(&myAccount)

	// 若 找不到帳號 則 回傳錯誤
	if database.GetAccountByEmail(myAccount.Email).ID != 0 {
		return context.SendStatus(fiber.StatusBadRequest)
	}

	// 若 新增成功 則 回傳成功
	if database.AddAccount(myAccount) {
		return context.SendStatus(fiber.StatusOK)
	}

	// 回傳錯誤
	return context.SendStatus(fiber.StatusInternalServerError)
}