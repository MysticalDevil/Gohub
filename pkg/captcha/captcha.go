// Package captcha Handle image verification code logic
package captcha

import (
	"github.com/mojocn/base64Captcha"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"sync"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// once Make sure the internalCaptcha object is only initialized once
var once sync.Once

// internalCaptcha Captcha object used internally
var internalCaptcha *Captcha

// NewCaptcha Singleton mode acquisition
func NewCaptcha() *Captcha {
	once.Do(func() {
		// Initialize Captcha object
		internalCaptcha = &Captcha{}

		// Use the global Redis object and configure the prefix for storing the Key
		store := RedisStore{
			RedisClient: redis.Redis,
			KeyPrefix:   config.GetString("app.name") + ":captcha:",
		}

		// Configure base64Captcha driver information
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("captcha.height"),
			config.GetInt("captcha.width"),
			config.GetInt("captcha.length"),
			config.GetFloat64("captcha.maxSkew"), // The maximum tilt angle of the numbers
			config.GetInt("captcha.dotCount"),    // The number of confusion points in the background of the image
		)

		// Instantiate base64Captcha and assign it to the internalCaptcha object used internally
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})
	return internalCaptcha
}

// GenerateCaptcha Generate image verification code
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

// VerifyCaptcha Verify that the image verification code is correct
func (c *Captcha) VerifyCaptcha(id, answer string) (match bool) {
	if !app.IsProduction() && id == config.GetString("captcha.testing_key") {
		return true
	}
	// The third parameter is whether to delete after verification
	return c.Base64Captcha.Verify(id, answer, false)
}
