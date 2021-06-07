//Package notify has functions and types used for sending notifications on different communication channel
package notify

import (
	"fmt"
	"net/smtp"
)

const (
	smtpServerAddress = "smtp.gmail.com"
	smtpServerPort    = "587"
)

type Email struct {
	ID   string
	Pass string
}

// NewEmail returns the instance of Email
func NewEmail(id, pass string) Notifier {
	return &Email{
		ID:   id,
		Pass: pass,
	}
}

// SendMessage takes message body and send it to the given email-id
func (e *Email) SendMessage(body string) error {
	msg := "From: " + e.ID + "\n" +
		"To: " + e.ID + "\n" +
		"Subject: Vaccination slots are available\n\n" +
		"Vaccination slots are available at the following centers:\n\n" +
		body

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpServerAddress, smtpServerPort),
		smtp.PlainAuth("", e.ID, e.Pass, smtpServerAddress),
		e.ID, []string{e.ID}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
