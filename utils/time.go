package utils

import "time"

// 取得當前時區
func GetCurrentTime() time.Time {
	now := time.Now()
	_, offset := now.Zone()
	return now.Add(time.Duration(offset) * time.Second)
}
