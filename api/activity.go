package api

import (
	"fmt"
	"time"

	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

var activityCache = cache.New(30*time.Minute, 60*time.Minute)

// 取得運動員所有活動
func GetActivities(context *fiber.Ctx) error {
	athleteID, _ := context.ParamsInt("athleteid") // 運動員主鍵
	cacheKey := fmt.Sprintf("activity-%d", athleteID)

	if cachedData, found := activityCache.Get(cacheKey); found {
		return context.JSON(cachedData)
	}

	myAthlete := database.GetAthleteByID(uint64(athleteID)) // 依運動員主鍵取得運動員

	_ = helper.FetchStravaActivities(myAthlete)

	activities := database.GetActivitiesByAthleteID(myAthlete.ID)

	activityCache.Set(cacheKey, activities, 30*time.Minute)

	return context.JSON(activities)
}

// 取得活動圈數
func GetActivityLaps(context *fiber.Ctx) error {
	athleteID, _ := context.ParamsInt("athleteid") // 運動員主鍵
	activityID := context.Params("activityid")     // 活動主鍵

	myAthlete := database.GetAthleteByID(uint64(athleteID)) // 依運動員主鍵取得運動員

	activity := database.GetActivityByID(activityID)
	laps := database.GetLaps(myAthlete.ID, activityID) // 取得跑步活動圈數
	// 若 沒有跑步圈數 則 爬取後再取得跑步活動圈數
	if len(laps) == 0 {
		_ = helper.FetchStravaLaps(myAthlete, activityID)
		laps = database.GetLaps(myAthlete.ID, activityID) // 取得跑步活動圈數
	}
	activity.Laps = laps

	return context.JSON(activity)
}
