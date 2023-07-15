package utils

import (
	"auth/helpers/constants"
	"bytes"
	"net/smtp"
)

func SendEmail(toEmails []string, subject string, body []byte) error {
	//sender data

	from := constants.Email
	password := constants.Password
	host := constants.Host
	port := constants.Port
	address := host + ":" + port
	var mail bytes.Buffer
	mail.WriteString(subject)
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n\n"
	mail.WriteString(mimeHeaders)
	mail.Write(body)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, toEmails, mail.Bytes())
	if err != nil {
		return err
	}
	return nil
}
