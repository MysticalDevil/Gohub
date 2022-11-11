// Package sms Send SMS
package sms

import (
	"gohub/pkg/config"
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
			Driver: &Aliyun{},
		}
	})

	return internalSMS
}

func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, config.GetStringMapString("sms.aliyun"))
}
