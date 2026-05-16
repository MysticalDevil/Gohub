// Package hash operation class
package hash

import (
	"gohub/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash Encrypt password with bcrypt
func BcryptHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogIf(err)

	return string(bytes)
}

// BcryptCheck Compare plaintext passwords to database hashes
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BcryptIsHashed(str string) bool {
	_, err := bcrypt.Cost([]byte(str))
	return err == nil
}
