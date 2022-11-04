package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
)

func PrepareEmailContent(wd string, subject string, toMail string) (MailContent, error) {
	// read kamp2022.html template
	t, err := template.ParseFiles(wd + "/templates/email.html")
	if err != nil {
		return MailContent{}, err
	}

	var body bytes.Buffer // buffer to write html to
	// execute template
	t.Execute(&body, struct {
		Key string
	}{
		Key: "",
	})

	// create dto
	mailContent := MailContent{
		ToName:      toMail,
		To:          toMail,
		Body:        body.Bytes(),
		Subject:     subject,
		MimeVersion: "1.0",
		ContentType: "text/html; charset=\"UTF-8\"",
	}

	return mailContent, nil
}

func SendEmailViaSmtp(mailConfig *MailConfig, mailContent MailContent) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		mailConfig.Username,
		mailConfig.Password,
		mailConfig.Host,
	)

	from := mail.Address{
		Name:    mailConfig.FromName,
		Address: mailConfig.FromMail,
	}
	to := mail.Address{
		Name:    mailContent.ToName,
		Address: mailContent.To,
	}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = mailContent.Subject
	header["MIME-Version"] = mailContent.MimeVersion
	header["Content-Type"] = mailContent.ContentType
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(mailContent.Body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		mailConfig.Host+":"+mailConfig.Port,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		return err
	}

	return nil
}
