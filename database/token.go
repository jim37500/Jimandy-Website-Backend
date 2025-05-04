package database

import "Jimandy-Website-Backend/model"

// 新增 token
func SaveTokens(token *model.Token) bool {
	return db.Save(token).Error == nil
}

// 依 access token 取得 token
func GetTokenByAccessToken(accessToken string) *model.Token {
	var token model.Token
	if db.Where("access_token = ?", accessToken).First(&token).Error != nil {
		return nil
	}
	return &token
}

// 依 refresh token 取得 token
func GetTokenByRefreshToken(refreshToken string) *model.Token {
	var token model.Token
	if db.Where("refresh_token = ?", refreshToken).First(&token).Error != nil {
		return nil
	}
	return &token
}

// 取得用戶的所有  token
func GetUserTokens(accountID uint) []model.Token {
	var tokens []model.Token
	db.Where("account_id = ? AND is_revoked = ?", accountID, false).Find(&tokens)
	return tokens
}

// 撤銷用戶的所有 token
func RevokeAllUserTokens(accountID uint) bool {
	// 撤銷所有 token
	if err := db.Model(&model.Token{}).Where("account_id = ?", accountID).Update("is_revoked", true).Error; err != nil {
		return false
	}

	return true
}
