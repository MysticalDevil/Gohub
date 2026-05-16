// Package user Store user model related logic
package user

import (
	"context"

	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/hash"
)

// User user model
type User struct {
	models.BaseModel

	Name string `json:"name,omitempty"`

	City         string `json:"city,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Avatar       string `json:"avatar,omitempty"`

	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}

func (userModel *User) Create(ctx context.Context) {
	database.DBWithContext(ctx).Create(&userModel)
}

// ComparePassword Is the password correct
func (userModel *User) ComparePassword(_password string) bool {
	return hash.BcryptCheck(_password, userModel.Password)
}

func (userModel *User) Save(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Save(&userModel)
	return result.RowsAffected
}

func (userModel *User) UpdateProfile(ctx context.Context, name, city, introduction string) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Model(&userModel).Updates(map[string]any{
		"name":         name,
		"city":         city,
		"introduction": introduction,
	})
	return result.RowsAffected
}

func (userModel *User) UpdateEmail(ctx context.Context, email string) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Model(&userModel).Update("email", email)
	return result.RowsAffected
}

func (userModel *User) UpdatePhone(ctx context.Context, phone string) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Model(&userModel).Update("phone", phone)
	return result.RowsAffected
}

func (userModel *User) UpdatePassword(ctx context.Context, password string) (rowsAffected int64) {
	hashedPassword := hash.BcryptHash(password)
	result := database.DBWithContext(ctx).Model(&userModel).Update("password", hashedPassword)
	return result.RowsAffected
}

func (userModel *User) UpdateAvatar(ctx context.Context, avatar string) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Model(&userModel).Update("avatar", avatar)
	return result.RowsAffected
}
