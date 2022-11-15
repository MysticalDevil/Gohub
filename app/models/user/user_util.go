package user

import "gohub/pkg/database"

// IsEmailExist Determine if Email has been registered
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email=?", email).Count(&count)
	return count > 0
}

// IsPhoneExist Determine if phone number has been registered
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone=?", phone).Count(&count)
	return count > 0
}

// GetByPhone Get users by phone number
func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone = ?", phone).First(&userModel)
	return
}

// GetByUtil Get users by phone/Email/username
func GetByUtil(loginID string) (userModel User) {
	database.DB.
		Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	return
}

// Get user by ID
func Get(idStr string) (userModel User) {
	database.DB.Where("id", idStr).First(&userModel)
	return
}
