package utils

import "time"

var taipeiLoc *time.Location

func init() {
	var err error
	taipeiLoc, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

// 取得當前時區
func GetCurrentTime() time.Time {
	return time.Now().In(taipeiLoc)
}
