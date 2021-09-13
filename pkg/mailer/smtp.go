package mailer

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"text/template"
)

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
}

func FromTemplate(filename string, data interface{}, subject string, from mail.Address, to []string) error {
	var message bytes.Buffer
	headers := make(map[string]string)

	config := SMTPConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
	}

	headers["Subject"] = subject
	headers["From"] = from.String()
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["MIME-Version"] = "1.0"

	for header, value := range headers {
		if _, err := fmt.Fprintf(&message, "%s: %s\r\n", header, value); err != nil {
			return err
		}
	}

	// Message body has to be separated by an additional line break
	message.WriteString("\r\n")
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	if err = tpl.Execute(&message, data); err != nil {
		return err
	}

	addr := net.JoinHostPort(config.Host, config.Port)
	auth := smtp.PlainAuth("", config.User, config.Pass, config.Host)

	return smtp.SendMail(addr, auth, config.User, to, message.Bytes())
}
