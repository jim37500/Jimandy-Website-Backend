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

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// 權杖預設逾時
const (
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

func setAccessToken(account *model.Account, deviceInfo string) string {
	now := time.Now()
	// 產生 access token
	accessClaims := jwt.MapClaims{
		"id":  account.ID,
		"exp": now.Add(accessTokenExpireDuration).Unix(),
		"jti": generateRandomString(16), // 加入隨機字串作為 JWT ID
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, _ := accessToken.SignedString(configuration.JWTKey)

	// 儲存 access token 到資料庫
	myAccessToken := model.AccessToken{
		AccountID:   account.ID,
		AccessToken: signedAccessToken,
		DeviceInfo:  deviceInfo,
		CreatedAt:   now,
		ExpiresAt:   now.Add(accessTokenExpireDuration),
	}

	if database.SaveAccessToken(&myAccessToken) {
		return signedAccessToken
	}
	return ""
}

func setRefreshToken(account *model.Account, deviceInfo string) string {
	now := time.Now()
	// 產生 refresh token
	refreshClaims := jwt.MapClaims{
		"id":  account.ID,
		"exp": now.Add(refreshTokenExpireDuration).Unix(),
		"jti": generateRandomString(16), // 加入隨機字串作為 JWT ID
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, _ := refreshToken.SignedString(configuration.JWTKey)

	// 儲存 refresh token 到資料庫
	myRefreshToken := model.RefreshToken{
		AccountID:    account.ID,
		RefreshToken: signedRefreshToken,
		DeviceInfo:   deviceInfo,
		CreatedAt:    now,
		ExpiresAt:    now.Add(refreshTokenExpireDuration),
	}
	if database.SaveRefreshToken(&myRefreshToken) {
		return signedRefreshToken
	}
	return ""
}

// 產生權杖
func GenerateTokens(account *model.Account, context *fiber.Ctx) fiber.Map {
	deviceInfo := getDeviceInfo(context)

	return fiber.Map{
		"accessToken":  setAccessToken(account, deviceInfo),
		"refreshToken": setRefreshToken(account, deviceInfo),
	}
}

// 刷新 access token
func RefreshToken(context *fiber.Ctx) error {
	refreshToken := GetTokenFromHeader(context)
	if refreshToken == "" {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 檢查 refresh token 是否有效
	dbToken := database.GetRefreshTokenByToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked || time.Now().After(dbToken.ExpiresAt) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 獲取用戶信息
	account := database.GetAccountByID(uint64(dbToken.AccountID))
	if account.ID == 0 {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 生成新的 access token
	signedAccessToken := setAccessToken(&account, dbToken.DeviceInfo)

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
	dbToken := database.GetRefreshTokenByToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 撤銷當前的 refresh token
	dbToken.IsRevoked = true
	database.SaveRefreshToken(dbToken)

	// 撤銷對應的 access token
	accessTokens := database.GetUserAccessTokens(dbToken.AccountID)
	for _, token := range accessTokens {
		if token.DeviceInfo == dbToken.DeviceInfo {
			token.IsRevoked = true
			database.SaveAccessToken(&token)
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
	dbToken := database.GetAccessTokenByToken(token)
	if dbToken == nil || dbToken.IsRevoked || time.Now().After(dbToken.ExpiresAt) {
		return false
	}

	return true
}

// 取得用戶的所有裝置
func GetUserDevices(context *fiber.Ctx) error {
	accountID := uint(context.Locals("id").(float64))
	devices := database.GetUserAccessTokens(accountID)

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
	dbToken := database.GetRefreshTokenByToken(refreshToken)
	if dbToken == nil || dbToken.IsRevoked || time.Now().After(dbToken.ExpiresAt) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	deviceInfo := getDeviceInfo(context)

	// 撤銷指定裝置的所有 access token
	tokens := database.GetUserAccessTokens(dbToken.AccountID)
	for _, token := range tokens {
		if token.DeviceInfo == deviceInfo {
			token.IsRevoked = true
			database.SaveAccessToken(&token)
		}
	}

	// 撤銷指定裝置的所有 refresh token
	refreshTokens := database.GetUserRefreshTokens(dbToken.AccountID)
	for _, token := range refreshTokens {
		if token.DeviceInfo == deviceInfo {
			token.IsRevoked = true
			database.SaveRefreshToken(&token)
		}
	}

	return context.SendStatus(fiber.StatusOK)
}
