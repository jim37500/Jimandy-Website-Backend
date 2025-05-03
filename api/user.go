package api

import (
	"Jimandy-Website-Backend/database"

	"github.com/gofiber/fiber/v2"
)

// 取得當前使用者
func GetCurrentUser(context *fiber.Ctx) error {
	id := uint64(context.Locals("id").(float64))

	account := database.GetAccountByID(id)

	return context.JSON(fiber.Map{
		"Token":   GenerateTokens(&account, context),
		"Account": account,
		"Athlete": database.GetAthleteByAccountID(uint(id)),
	})
}
