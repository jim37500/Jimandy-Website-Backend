package router

import (
	"path/filepath"

	"Jimandy-Website-Backend/api"
	"Jimandy-Website-Backend/configuration"
	"Jimandy-Website-Backend/data"
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
)

const (
	bodyLimit = 1 * 1024 * 1024 // 上傳大小限制1MB
)

var HttpApplication *fiber.App

// 監聽網頁服務
func Run() {
	HttpApplication = fiber.New(fiber.Config{BodyLimit: bodyLimit, UnescapePath: true, ErrorHandler: notFoundHandler})

	HttpApplication.Use(compress.New()) // 啟用壓縮
	HttpApplication.Use(recover.New())  // 啟用錯誤處理
	HttpApplication.Use(logger.New())   // 啟用日誌
	HttpApplication.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,https://jimandy-growth.com/",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 將靜態頁面綁定根目錄
	HttpApplication.Static("/", filepath.Join(configuration.ExecutPath, "dist"))

	setupRoute() // 設定路由

	_ = HttpApplication.Listen(":61018")
}

func bindAuthorized() {
	// 要求授權
	HttpApplication.Use(jwtware.New(jwtware.Config{
		SigningKey:     configuration.JWTKey,
		SuccessHandler: jwtSuccessHandler, // 權杖驗證成功後檢查授權
	}))
	apiGroup := HttpApplication.Group("")

	apiGroup.Get("/api/currentuser", api.GetCurrentUser)   // 取得當前使用者資訊
	apiGroup.Get("/api/devices", api.GetUserDevices)       // 取得用戶的所有裝置
	apiGroup.Post("/api/devices/logout", api.LogoutDevice) // 登出特定裝置

	// 綁定授權列表
	for _, acl := range AccessControlLists {
		apiGroup.Add(acl.Method, acl.Path, acl.Handler)
	}
}

// 未找到路徑轉到首頁
func notFoundHandler(context *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// 若 未找到路徑 則 回到首頁
	if code == fiber.StatusNotFound {
		return context.Redirect("/")
	}

	return context.Status(code).SendString(err.Error())
}

// 設定路由
func setupRoute() {
	HttpApplication.Post("/api/login", api.Login)          // 取得帳號權杖
	HttpApplication.Post("/api/refresh", api.RefreshToken) // 刷新 access token
	HttpApplication.Post("/api/logout", api.Logout)        // 登出

	bindAuthorized() // 綁定授權
}

// AccessControlLists 存取控制列表
var AccessControlLists = []data.AccessControlList{
	AccessControlListFactory("/api/athlete", fiber.MethodPost, "", api.AddAthlete),                                    // 新增運動員
	AccessControlListFactory("/api/activities/:athleteid", fiber.MethodGet, "", api.GetActivities),                    // 取得所有活動紀錄
	AccessControlListFactory("/api/activities/laps/:athleteid/:activityid", fiber.MethodGet, "", api.GetActivityLaps), // 取得活動紀錄圈數
}

// 存取控制列表 工廠方法
func AccessControlListFactory(path, method string, authorization string, handler fiber.Handler) data.AccessControlList {
	return data.AccessControlList{Path: path, Method: method, Authorization: authorization, Handler: handler}
}

// 權杖驗證成功後檢查授權
func jwtSuccessHandler(context *fiber.Ctx) error {
	token := api.GetTokenFromHeader(context)
	if !api.IsTokenValid(token) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	// 從資料庫獲取 token 信息
	dbToken := database.GetTokenByAccessToken(token)
	if dbToken == nil || dbToken.IsRevoked || utils.GetCurrentTime().After(dbToken.ExpiresAt) {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

	context.Locals("id", float64(dbToken.AccountID)) // 登入帳號

	return context.Next()
}
