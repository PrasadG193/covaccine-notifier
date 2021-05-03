package main

import (
	"net/smtp"
)

func sendMail(id, pass, body string) error {
	msg := "From: " + id + "\n" +
		"To: " + id + "\n" +
		"Subject: Vaccination slots are available\n\n" +
		"Vaccination slots are available at the following centers:\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", id, pass, "smtp.gmail.com"),
		id, []string{id}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
