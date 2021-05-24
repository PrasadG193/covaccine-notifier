package notify

import (
	"net/smtp"
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

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", e.Id, e.Pass, "smtp.gmail.com"),
		e.Id, []string{e.Id}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
