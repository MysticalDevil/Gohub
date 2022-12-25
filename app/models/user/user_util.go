package user

import (
	"github.com/gin-gonic/gin"
	"gohub/pkg/app"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
)

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

// GetByEmail Get users by email
func GetByEmail(email string) (userModel User) {
	database.DB.Where("email = ?", email).First(&userModel)
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

// All user data
func All() (users []User) {
	database.DB.Find(&users)
	return
}

// Paginate Pagination content
func Paginate(c *gin.Context, perPage int) (users []User, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(User{}),
		&users,
		app.V1URL(database.TableName(&User{})),
		perPage,
	)
	return
}
