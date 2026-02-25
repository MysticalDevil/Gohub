// Package mail Send e-mail
package mail

import (
	"sync"

	"gohub/pkg/config"
)

type From struct {
	Address string
	Name    string
}

type Email struct {
	From    From
	To      []string // Recipient
	Bcc     []string // Blind Carbon Copy
	Cc      []string // Carbon Copy
	Subject string
	Text    []byte // Plaintext message (optional)
	HTML    []byte // Html message (optional)
}

type Mailer struct {
	Driver Driver
}

var (
	once           sync.Once
	internalMailer *Mailer
)

// NewMailer Singleton mode acquisition
func NewMailer() *Mailer {
	once.Do(func() {
		internalMailer = &Mailer{
			Driver: &SMTP{},
		}
	})

	return internalMailer
}

func (mailer *Mailer) Send(email Email) bool {
	return mailer.Driver.Send(email, config.GetStringMapString("mail.smtp"))
}
