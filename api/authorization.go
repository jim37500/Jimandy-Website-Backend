package api

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"Jimandy-Website-Backend/configuration"
	"Jimandy-Website-Backend/data"
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/model"
	"Jimandy-Website-Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var (
	// 權杖預設逾時
	accessTokenExpireDuration  = time.Hour * 24 * 7      // 1小時
	refreshTokenExpireDuration = time.Hour * 24 * 7 * 30 // 一個月
)

// 產生隨機字串
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 取得裝置資訊
func getDeviceInfo(context *fiber.Ctx) string {
	userAgent := context.Get("User-Agent")
	ip := context.IP()
	return userAgent + "|" + ip
}

// 從 Authorization header 取得 token
func GetTokenFromHeader(context *fiber.Ctx) string {
	token := context.Get("Authorization")
	if token == "" {
		return ""
	}
	// 移除 Bearer 前綴
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}
	return token
}

func Login(context *fiber.Ctx) error {
	var myLoginData data.Login
	_ = context.BodyParser(&myLoginData)

	// 若 沒Email 則 回錯誤
	if myLoginData.Email == "" {
		return context.SendStatus(fiber.StatusBadRequest)
	}

	account := database.GetAccountByEmail(myLoginData.Email) // 依 登入帳號 取得 帳號
	if account.ID == 0 {
		account = model.Account{
			Name:  myLoginData.Name,
			Email: myLoginData.Email,
		}
		if !database.AddAccount(account) {
			return context.SendStatus(fiber.StatusBadRequest)
		}
	}

	return context.JSON(GenerateTokens(&account, context))
}

func setAccessToken(account *model.Account) string {
	now := utils.GetCurrentTime()
	// 產生 access token
	accessClaims := jwt.MapClaims{
		"id":  account.ID,
		"exp": now.Add(accessTokenExpireDuration).Unix(),
		"jti": generateRandomString(16), // 加入隨機字串作為 JWT ID
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, _ := accessToken.SignedString(configuration.JWTKey)

	return signedAccessToken
}

func setRefreshToken(account *model.Account) string {
	now := utils.GetCurrentTime()
	// 產生 refresh token
	refreshClaims := jwt.MapClaims{
		"id":  account.ID,
		"exp": now.Add(refreshTokenExpireDuration).Unix(),
		"jti": generateRandomString(16), // 加入隨機字串作為 JWT ID
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, _ := refreshToken.SignedString(configuration.JWTKey)

	return signedRefreshToken
}

// 產生權杖
func GenerateTokens(account *model.Account, context *fiber.Ctx) fiber.Map {
	deviceInfo := getDeviceInfo(context)
	now := utils.GetCurrentTime()

	// 產生 tokens
	accessToken := setAccessToken(account)
	refreshToken := setRefreshToken(account)

	// 儲存 tokens 到資料庫
	myToken := model.Token{
		AccountID:    account.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceInfo:   deviceInfo,
		CreatedAt:    now,
		ExpiresAt:    now.Add(accessTokenExpireDuration),
	}

	if database.SaveTokens(&myToken) {
		return fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}
	}
	return fiber.Map{}
}

// 刷新 access token
func RefreshToken(context *fiber.Ctx) error {
	refreshToken := GetTokenFromHeader(context)
	if refreshToken == "" {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 檢查 refresh token 是否有效
	dbToken := database.GetTokenByRefreshToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked || utils.GetCurrentTime().After(dbToken.ExpiresAt) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 獲取用戶信息
	account := database.GetAccountByID(uint64(dbToken.AccountID))
	if account.ID == 0 {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 生成新的 access token
	signedAccessToken := setAccessToken(&account)
	dbToken.AccessToken = signedAccessToken
	database.SaveTokens(dbToken)

	return context.JSON(fiber.Map{
		"accessToken": signedAccessToken,
	})
}

// 登出
func Logout(context *fiber.Ctx) error {
	refreshToken := GetTokenFromHeader(context)
	if refreshToken == "" {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 檢查 refresh token 是否有效
	dbToken := database.GetTokenByRefreshToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 撤銷當前的 refresh token
	dbToken.IsRevoked = true
	database.SaveTokens(dbToken)

	// 撤銷對應的 access token
	accessTokens := database.GetUserTokens(dbToken.AccountID)
	for _, token := range accessTokens {
		if token.DeviceInfo == dbToken.DeviceInfo {
			token.IsRevoked = true
			database.SaveTokens(&token)
		}
	}

	return context.SendStatus(fiber.StatusOK)
}

// 檢查 token 是否有效
func IsTokenValid(token string) bool {
	if token == "" {
		return false
	}

	// 檢查 token 是否在資料庫中且未撤銷
	dbToken := database.GetTokenByAccessToken(token)
	if dbToken == nil || dbToken.IsRevoked || time.Now().After(dbToken.ExpiresAt) {
		return false
	}

	return true
}

// 取得用戶的所有裝置
func GetUserDevices(context *fiber.Ctx) error {
	accountID := uint(context.Locals("id").(float64))
	devices := database.GetUserTokens(accountID)

	var deviceList []fiber.Map
	for _, token := range devices {
		deviceInfo := strings.Split(token.DeviceInfo, "|")
		deviceList = append(deviceList, fiber.Map{
			"deviceInfo": deviceInfo[0],
			"ip":         deviceInfo[1],
			"lastLogin":  token.CreatedAt,
		})
	}

	return context.JSON(deviceList)
}

// 登出特定裝置
func LogoutDevice(context *fiber.Ctx) error {
	refreshToken := GetTokenFromHeader(context)
	if refreshToken == "" {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 檢查 refresh token 是否有效
	dbToken := database.GetTokenByRefreshToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked || time.Now().After(dbToken.ExpiresAt) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	deviceInfo := getDeviceInfo(context)

	// 撤銷指定裝置的所有 token
	tokens := database.GetUserTokens(dbToken.AccountID)
	for _, token := range tokens {
		if token.DeviceInfo == deviceInfo {
			token.IsRevoked = true
			database.SaveTokens(&token)
		}
	}

	return context.SendStatus(fiber.StatusOK)
}
