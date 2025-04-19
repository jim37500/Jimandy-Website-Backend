package api

import (
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/model"

	"github.com/gofiber/fiber/v2"
)

func GoogleLogin(context *fiber.Ctx) error {
	var myAccount model.Account
	_ = context.BodyParser(&myAccount)

	// 若 沒Email 則 回錯誤
	if myAccount.Email == "" {
		return context.SendStatus(fiber.StatusBadRequest)
	}

	account := database.GetAccountByEmail(myAccount.Email)

	if account.Name == "" {
		if !database.AddAccount(myAccount) {
			return context.SendStatus(fiber.StatusBadRequest)
		}
		account = database.GetAccountByEmail(myAccount.Email)
	}

	return context.JSON(GenerateToken(&account))
}