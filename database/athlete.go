package database

import "Jimandy-Website-Backend/model"

// 依主鍵取得運動員
func GetAthleteByID(athleteID uint64) (athlete model.Athlete) {
	db.First(&athlete, athleteID)

	return
}

func GetAthleteByAccountID(accountID uint) (athlete model.Athlete) {
	db.Where("account_id = ?", accountID).First(&athlete)

	return
}

func AddAthlete(myAthlete *model.Athlete) bool {
	return db.Save(myAthlete).Error == nil
}

// 更新運動員權杖
func UpdateAthleteToken(athleteID uint64, accessToken string, refreshToken string) bool {
	var athlete model.Athlete
	db.First(&athlete, athleteID)

	if athlete.ID == 0 {
		return false
	}

	athlete.AccessToken = accessToken
	athlete.RefreshToken = refreshToken

	return db.Save(&athlete).Error == nil
}