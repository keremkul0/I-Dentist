package init_func

import (
	"dental-clinic-system/config"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type goMail struct {
	dialer *gomail.Dialer
	from   string
}

func SetupMailDialer(emailConfig config.EmailConfig) *goMail {
	mailDialer := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.User, emailConfig.Password)
	if mailDialer == nil {
		log.Fatal().Msg("Error setting up mail dialer")
		panic("Error setting up mail dialer")
	}

	// Test connection
	conn, err := mailDialer.Dial()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the mail server")
		panic(fmt.Sprintf("Failed to connect to the mail server: %v", err))
	}
	defer func(conn gomail.SendCloser) {
		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close the connection")
		}
	}(conn)

	log.Info().Msg("Mail dialer setup successful and connection test passed")

	return &goMail{dialer: mailDialer, from: emailConfig.User}
}

func (g *goMail) SendMail(to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	// Send the email
	if err := g.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
