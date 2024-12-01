package emailsender

import (
	"net/smtp"

	"auth-service/config"
)

type EmailSender func(emailAddress, message string) error

func NewEmailSender(cfg *config.Config) EmailSender {
	return func(emailAddress, message string) error {
		from := cfg.EmailAdress
		password := cfg.EmailPassword

		to := []string{emailAddress}

		host := cfg.SMTPHost
		port := cfg.SMTPPort
		address := host + ":" + port

		auth := smtp.PlainAuth("", from, password, host)

		err := smtp.SendMail(address, auth, from, to, []byte(message))
		if err != nil {
			return err
		}

		return nil
	}
}
