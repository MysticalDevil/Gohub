package user

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
	"gorm.io/gorm"
)

func IsEmailExist(ctx context.Context, email string) (bool, error) {
	return models.Exists[User](ctx, "email", email)
}

func IsPhoneExist(ctx context.Context, phone string) (bool, error) {
	return models.Exists[User](ctx, "phone", phone)
}

func GetByPhone(ctx context.Context, phone string) (userModel User, err error) {
	return models.GetBy[User](ctx, "phone", phone)
}

func GetByEmail(ctx context.Context, email string) (userModel User, err error) {
	return models.GetBy[User](ctx, "email", email)
}

func GetByUtil(ctx context.Context, loginID string) (userModel User, err error) {
	result := database.DBWithContext(ctx).
		Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return userModel, nil
		}
		return userModel, result.Error
	}
	return userModel, nil
}

func Get(ctx context.Context, idStr string) (userModel User, err error) {
	return models.Get[User](ctx, idStr)
}

func All(ctx context.Context) (users []User, err error) {
	return models.All[User](ctx)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (users []User, paging paginator.Paging) {
	return models.Paginate[User](ctx, c, limit)
}
