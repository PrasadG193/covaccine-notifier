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
	Id   string
	Pass string
}

func NewEmail(id, pass string) Notifier {
	return &Email{
		Id:   id,
		Pass: pass,
	}
}

func (e *Email) SendMessage(body string) error {
	msg := "From: " + e.Id + "\n" +
		"To: " + e.Id + "\n" +
		"Subject: Vaccination slots are available\n\n" +
		"Vaccination slots are available at the following centers:\n\n" +
		body

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpServerAddress, smtpServerPort),
		smtp.PlainAuth("", e.Id, e.Pass, smtpServerAddress),
		e.Id, []string{e.Id}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
