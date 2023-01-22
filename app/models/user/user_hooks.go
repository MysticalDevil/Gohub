package user

import (
	"gohub/pkg/hash"
	"gorm.io/gorm"
)

// BeforeSave Gorm's model hook, called before creating and updating models
func (userModel *User) BeforeSave(_ *gorm.DB) (err error) {
	if !hash.BcryptIsHashed(userModel.Password) {
		userModel.Password = hash.BcryptHash(userModel.Password)
	}

	return
}
