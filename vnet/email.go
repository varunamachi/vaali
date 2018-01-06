package vnet

import (
	"fmt"
	"net/smtp"

	"github.com/varunamachi/vaali/vlog"
)

//SendEmail - sends an email with given information. Uses the package level
//variable emainConfig for SMTP configuration - smtp.gmail.com:587
func SendEmail(to, subject, meesage string) (err error) {
	msg := "From: " + emailConfig.From + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
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
