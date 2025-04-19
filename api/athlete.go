package api

import (
	"Jimandy-Website-Backend/helper"
	"Jimandy-Website-Backend/model"

	"github.com/gofiber/fiber/v2"
)

func AddAthlete(context *fiber.Ctx) error {
	accountID := context.Locals("id").(float64) // 登入帳號主鍵

	var myAthlete model.Athlete
	_ = context.BodyParser(&myAthlete)

	myAthlete.AccountID = uint(accountID)

	// 若 授權Strava成功 則 回傳成功
	if helper.AuthorizeStrava(myAthlete) {
		return context.SendStatus(fiber.StatusOK)
	}

	// 否則 回傳錯誤
	return context.SendStatus(fiber.StatusInternalServerError)
}
