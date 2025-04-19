package model

import (
	"time"
)

// 跑步活動
type Activity struct {
	ID                 uint64    `gorm:"primarykey"`
	AthleteID          uint64    `gorm:"comment:運動員主鍵"`
	Name               string    `gorm:"comment:活動名稱"`
	SportType          string    `gorm:"comment:活動類別"`
	Date               time.Time `gorm:"comment:時間"`
	ElapsedTime        int       `gorm:"comment:持續時間(秒)"`
	MovingTime         int       `gorm:"comment:運動時間(秒)"`
	Distance           int       `gorm:"comment:距離(m)"`
	TotalElevationGain int       `gorm:"comment:總爬升(m)"`
	AverageSpeed       float32   `gorm:"comment:平均速度(km/h)"`
	MaxSpeed           float32   `gorm:"comment:最大速度(km/h)"`
	AverageCadence     int       `gorm:"comment:踏頻(spm)"`
	AverageHeartrate   int       `gorm:"comment:平均心率(bpm)"`
	MaxHeartrate       int       `gorm:"comment:最大心率(bpm)"`
	AverageWatts       int       `gorm:"comment:平均功率(bpm)"`
	MaxWatts           int       `gorm:"comment:最大功率(bpm)"`
	AverageTemperature int       `gorm:"comment:平均溫度(°C)"`
	Polyline           string    `gorm:"comment:路線"`
	Laps               []Lap     `gorm:"foreignKey:ActivityID"`
}

// 跑步圈數
type Lap struct {
	ID                 uint64  `gorm:"primarykey"`
	ActivityID         uint64  `gorm:"comment:活動主鍵"`
	AthleteID          uint64  `gorm:"comment:運動員主鍵"`
	ElapsedTime        int     `gorm:"comment:持續時間(秒)"`
	MovingTime         int     `gorm:"comment:運動時間(秒)"`
	Distance           int     `gorm:"comment:距離(m)"`
	AverageSpeed       float32 `gorm:"comment:平均速度(km/h)"`
	MaxSpeed           float32 `gorm:"comment:最大速度(km/h)"`
	AverageCadence     int     `gorm:"comment:踏頻(spm)"`
	AverageHeartrate   int     `gorm:"comment:平均心率(bpm)"`
	MaxHeartrate       int     `gorm:"comment:最大心率(bpm)"`
	AverageWatts       int     `gorm:"comment:平均功率(bpm)"`
	TotalElevationGain int     `gorm:"comment:爬升(cm)"`
}
