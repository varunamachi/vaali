package vnet

import (
	"fmt"
	"net/smtp"

	"github.com/varunamachi/vaali/vlog"
)

//"smtp.gmail.com:587"

func sendEmail(to string, meesage string) (err error) {
	msg := "From: " + emailConfig.From + "\n" +
		"To: " + to + "\n" +
		"Subject: Sparrow Registration\n\n" +
		meesage
	smtpURL := fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort)
	auth := smtp.PlainAuth("",
		emailConfig.From,
		emailConfig.Password,
		emailConfig.SMTPHost)
	err = smtp.SendMail(
		smtpURL,
		auth,
		emailConfig.From,
		[]string{to},
		[]byte(msg))
	return vlog.LogError("Net:EMail", err)
}
