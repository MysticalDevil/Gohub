package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
)

// IsEmailExist Determine if Email has been registered
func IsEmailExist(ctx context.Context, email string) bool {
	return models.Exists[User](ctx, "email", email)
}

// IsPhoneExist Determine if phone number has been registered
func IsPhoneExist(ctx context.Context, phone string) bool {
	return models.Exists[User](ctx, "phone", phone)
}

// GetByPhone Get users by phone number
func GetByPhone(ctx context.Context, phone string) (userModel User) {
	return models.GetBy[User](ctx, "phone", phone)
}

// GetByEmail Get users by email
func GetByEmail(ctx context.Context, email string) (userModel User) {
	return models.GetBy[User](ctx, "email", email)
}

// GetByUtil Get users by phone/Email/username
func GetByUtil(ctx context.Context, loginID string) (userModel User) {
	database.DBWithContext(ctx).
		Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	return
}

// Get user by ID
func Get(ctx context.Context, idStr string) (userModel User) {
	return models.Get[User](ctx, idStr)
}

// All user data
func All(ctx context.Context) (users []User) {
	return models.All[User](ctx)
}

// Paginate Pagination content
func Paginate(ctx context.Context, c *gin.Context, limit int) (users []User, paging paginator.Paging) {
	return models.Paginate[User](ctx, c, limit)
}
