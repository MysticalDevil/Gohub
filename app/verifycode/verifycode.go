// Package verifycode Used to send mobile phone verification code and email verification code
package verifycode

import (
	"fmt"
	"strings"
	"sync"

	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/mail"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
)

type VerifyCode struct {
	Store Store
}

var (
	once               sync.Once
	internalVerifyCode *VerifyCode
)

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
		Template: "",
		Data:     map[string]string{"code": code},
		Content:  fmt.Sprintf("Your verification code is %v", code),
	})
}

// SendEmail Send Email verification code, example:
//
//	verifycode.NewVerifyCode().SendEmail(request.Email)
func (vc *VerifyCode) SendEmail(email string) error {
	// generate verify code
	code := vc.generateVerifyCode(email)

	if !app.IsProduction() && strings.HasSuffix(email, config.GetString("verifycode.debug_email_suffix")) {
		return nil
	}
	content := fmt.Sprintf("<h1>Your email verification code is %v </h1>", code)
	// Send email
	mail.NewMailer().Send(mail.Email{
		From: mail.From{
			Address: config.GetString("mail.from.address"),
			Name:    config.GetString("mail.from.name"),
		},
		To:      []string{email},
		Subject: "Email verify code",
		HTML:    []byte(content),
	})

	return nil
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
