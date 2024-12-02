package emailsender

import (
	"auth-service/config"

	"github.com/sirupsen/logrus"
)

type EmailSender func(emailAddress, message string) error

func NewEmailSender(cfg *config.Config) EmailSender {
	// return func(emailAddress, message string) error {
	// 	from := cfg.EmailAdress
	// 	password := cfg.EmailPassword

	// 	to := []string{emailAddress}

	// 	host := cfg.SMTPHost
	// 	port := cfg.SMTPPort
	// 	address := host + ":" + port

	// 	auth := smtp.PlainAuth("", from, password, host)

	// 	err := smtp.SendMail(address, auth, from, to, []byte(message))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	// }

	return func(emailAddress, message string) error {
		logrus.Infof("Sending email warning to %s", emailAddress)
		return nil
	}
}
