package user

import "gohub/pkg/database"

// IsEmailExist Determine if Email has been registered
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("emil=?", email).Count(&count)
	return count > 0
}

// IsPhoneExist Determine if phone number has been registered
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone=?", phone).Count(&count)
	return count > 0
}
