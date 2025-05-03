package database

import "Jimandy-Website-Backend/model"

// 儲存 access token
func SaveAccessToken(token *model.AccessToken) bool {
	return db.Save(token).Error == nil
}

// 儲存 refresh token
func SaveRefreshToken(token *model.RefreshToken) bool {
	return db.Save(token).Error == nil
}

// 依 access token 取得 token
func GetAccessTokenByToken(accessToken string) *model.AccessToken {
	var token model.AccessToken
	if db.Where("access_token = ?", accessToken).First(&token).Error != nil {
		return nil
	}
	return &token
}

// 依 refresh token 取得 token
func GetRefreshTokenByToken(refreshToken string) *model.RefreshToken {
	var token model.RefreshToken
	if db.Where("refresh_token = ?", refreshToken).First(&token).Error != nil {
		return nil
	}
	return &token
}

// 取得用戶的所有 access token
func GetUserAccessTokens(accountID uint) []model.AccessToken {
	var tokens []model.AccessToken
	db.Where("account_id = ? AND is_revoked = ?", accountID, false).Find(&tokens)
	return tokens
}

// 取得用戶的所有 refresh token
func GetUserRefreshTokens(accountID uint) []model.RefreshToken {
	var tokens []model.RefreshToken
	db.Where("account_id = ? AND is_revoked = ?", accountID, false).Find(&tokens)
	return tokens
}

// 撤銷用戶的所有 token
func RevokeAllUserTokens(accountID uint) bool {
	// 撤銷所有 access token
	if err := db.Model(&model.AccessToken{}).Where("account_id = ?", accountID).Update("is_revoked", true).Error; err != nil {
		return false
	}
	// 撤銷所有 refresh token
	if err := db.Model(&model.RefreshToken{}).Where("account_id = ?", accountID).Update("is_revoked", true).Error; err != nil {
		return false
	}
	return true
}
