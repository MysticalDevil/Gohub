package mail

type Driver interface {
	// Send Check verify code
	Send(email Email, config map[string]string) bool
}
