package database

import (
	"math"
	"time"

	"Jimandy-Website-Backend/data"
	"Jimandy-Website-Backend/model"
)

// 取得運動員最近一次的活動
func GetLatestActivityDate(athleteID uint64) (date time.Time) {
	db.Model(&model.Activity{}).Select("date").Where("athlete_id = ?", athleteID).Order("date DESC").First(&date)

	return
}

// 依主鍵取得活動
func GetActivityByID(activityID string) (activity model.Activity) {
	db.First(&activity, activityID)

	return
}

// 依運動員主鍵取得活動
func GetActivitiesByAthleteID(athleteID uint64) (activities []model.Activity) {
	db.Where("athlete_id = ?", athleteID).Order("date DESC").Find(&activities)

	return
}

// 取得跑步活動圈數
func GetLaps(athleteID uint64, activityID string) (runLaps []model.Lap) {
	db.Where("athlete_id = ? AND activity_id = ?", athleteID, activityID).Find(&runLaps)

	return
}

// 新增活動
func AddActivity(activities []data.Activities) bool {
	var myActivities []model.Activity
	result := true

	for _, actitvity := range activities {
		date, _ := time.Parse(time.RFC3339, actitvity.Date)

		myActivities = append(myActivities, model.Activity{
			ID:                 actitvity.ActivityID,
			AthleteID:          actitvity.Athlete.AthleteID,
			Name:               actitvity.Name,
			SportType:          actitvity.SportType,
			Date:               date,
			ElapsedTime:        actitvity.ElapsedTime,
			MovingTime:         actitvity.MovingTime,
			Distance:           int(math.Round(float64(actitvity.Distance))),
			TotalElevationGain: int(math.Round(float64(actitvity.TotalElevationGain))),
			AverageSpeed:       actitvity.AverageSpeed * 3.6,
			MaxSpeed:           actitvity.MaxSpeed * 3.6,
			AverageCadence:     int(math.Round(float64(actitvity.AverageCadence * 2))),
			AverageHeartrate:   int(math.Round(float64(actitvity.AverageHeartrate))),
			MaxHeartrate:       int(math.Round(float64(actitvity.MaxHeartrate))),
			AverageWatts:       int(math.Round(float64(actitvity.AverageWatts))),
			MaxWatts:           int(math.Round(float64(actitvity.MaxWatts))),
			AverageTemperature: actitvity.AverageTemperature,
			Polyline:           actitvity.RouteMap.Polyline,
		})
	}

	if len(myActivities) > 0 && db.Create(&myActivities).Error != nil {
		result = false
	}

	return result
}

// 新增圈數
func AddLaps(laps []data.Lap) bool {
	var myRunLaps []model.Lap
	for _, runLap := range laps {
		myRunLaps = append(myRunLaps, model.Lap{
			ID:                 runLap.LapID,
			ActivityID:         runLap.Activity.ActivityID,
			AthleteID:          runLap.Athlete.AthleteID,
			ElapsedTime:        runLap.ElapsedTime,
			MovingTime:         runLap.MovingTime,
			Distance:           int(math.Round(float64(runLap.Distance))),
			AverageSpeed:       runLap.AverageSpeed * 3.6,
			MaxSpeed:           runLap.MaxSpeed * 3.6,
			AverageCadence:     int(math.Round(float64(runLap.AverageCadence * 2))),
			AverageHeartrate:   int(math.Round(float64(runLap.AverageHeartrate))),
			MaxHeartrate:       int(math.Round(float64(runLap.MaxHeartrate))),
			AverageWatts:       int(math.Round(float64(runLap.AverageWatts))),
			TotalElevationGain: int(math.Round(float64(runLap.TotalElevationGain * 10))),
		})
	}
	return db.Create(&myRunLaps).Error == nil
}
