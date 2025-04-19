package api

import (
	"time"

	"Jimandy-Website-Backend/configuration"
	"Jimandy-Website-Backend/data"
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// 權杖預設逾時
const tokenExpireDuration = time.Hour * 24 * 7 // 一周

func GetToken(context *fiber.Ctx) error {
	var myLoginData data.Login
	_ = context.BodyParser(&myLoginData)

	// 若 沒Email 則 回錯誤
	if myLoginData.Email == "" {
		return context.SendStatus(fiber.StatusBadRequest)
	}

	account := database.GetAccountByEmail(myLoginData.Email) // 依 登入帳號 取得 帳號

	// 若 沒有密碼 或 密碼不正確 或 帳號未啟用 則 回錯誤
	// if len(account.Password) == 0 || account.Password != myLoginData.Password || account.Status != model.Enabled {
	if len(account.Password) == 0 || account.Password != myLoginData.Password {
		return context.SendStatus(fiber.StatusBadRequest)
	}

	return context.JSON(GenerateToken(&account))
}

// 產生權杖
func GenerateToken(account *model.Account) fiber.Map {
	claims := jwt.MapClaims{
		"id":  account.ID, // 登入帳號主鍵
		"exp": time.Now().Add(tokenExpireDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(configuration.JWTKey)

	return fiber.Map{"Token": signedToken}
}
