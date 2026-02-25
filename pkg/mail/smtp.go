package mail

import (
	"fmt"
	"net/smtp"

	emailPKG "github.com/jordan-wright/email"
	"gohub/pkg/logger"
)

// SMTP Implement the email.Driver interface
type SMTP struct{}

// Send Implement the Send method of the email.Driver interface
func (s *SMTP) Send(email Email, config map[string]string) bool {
	e := emailPKG.NewEmail()

	e.From = fmt.Sprintf("%v <%v>", email.From.Name, email.From.Address)
	e.To = email.To
	e.Bcc = email.Bcc
	e.Cc = email.Cc
	e.Subject = email.Subject
	e.Text = email.Text
	e.HTML = email.HTML

	logger.DebugJSON("Send Email", "Send email details", e)

	err := e.Send(
		fmt.Sprintf("%v:%v", config["host"], config["port"]),
		smtp.PlainAuth(
			"",
			config["username"],
			config["password"],
			config["host"],
		),
	)
	if err != nil {
		logger.ErrorString("Send Email", "Error sending email", err.Error())
		return false
	}

	logger.DebugString("Send Email", "Send email success", "")
	return true
}
