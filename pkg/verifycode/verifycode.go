// Package verifycode Used to send mobile phone verification code and email verification code
package verifycode

import (
	"gohub/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
	"strings"
	"sync"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode Singleton mode acquisition
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				KeyPrefix:   config.GetString("app.name") + ":verifycode:",
			},
		}
	})

	return internalVerifyCode
}

// SendSMS Send SMS verification code, example:
//
//	verifycode.NewVerifyCode().SendSMS(request.Phone)
func (vc *VerifyCode) SendSMS(phone string) bool {
	// Generate verification code
	code := vc.generateVerifyCode(phone)

	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	// Send sms
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
}

// CheckAnswer Check whether the verification code submitted by the user is correct
func (vc *VerifyCode) CheckAnswer(key, answer string) bool {
	logger.DebugJSON("Verify Code", "Check verify code", map[string]string{key: answer})

	if !app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix")) ||
			strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}

	return vc.Store.Verify(key, answer, false)
}

// generateVerifyCode Generate verify code, and store in Redis
func (vc *VerifyCode) generateVerifyCode(key string) string {
	// generate random code
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("Verify Code", "Generate verify code", map[string]string{key: code})

	// Store the verification code and KEY in Redis and set the expiration time
	vc.Store.Set(key, code)
	return code
}
