package sms

type Driver interface {
	// Send SMS
	Send(phone string, message Message, config map[string]string) bool
}
