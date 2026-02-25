// Package sms Send SMS
package sms

import (
	"gohub/pkg/logger"
	"sync"
)

// Message SMS struct
type Message struct {
	Template string
	Data     map[string]string

	Content string
}

// SMS Action class for sending SMS
type SMS struct {
	Driver Driver
}

// once singleton pattern
var once sync.Once

// internalSMS SMS object used internally
var internalSMS *SMS

// NewSMS Singleton mode acquisition
func NewSMS() *SMS {
	once.Do(func() {
		internalSMS = &SMS{
			Driver: &Noop{},
		}
	})

	return internalSMS
}

func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, map[string]string{})
}

// Noop driver disables SMS sending when no provider is configured.
type Noop struct{}

func (noop *Noop) Send(phone string, message Message, _ map[string]string) bool {
	logger.WarnString("SMS", "Send", "SMS provider disabled; message not sent")
	return false
}
