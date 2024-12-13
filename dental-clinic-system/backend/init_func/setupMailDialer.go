package init_func

import (
	"dental-clinic-system/config"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

func SetupMailDialer(emailConfig *config.EmailConfig) *gomail.Dialer {
	mailDialer := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.User, emailConfig.Password)
	if mailDialer == nil {
		log.Fatal().Msg("Error setting up mail dialer")
		panic("Error setting up mail dialer")
	}
	log.Info().Msg("Mail dialer setup successful")
	return mailDialer
}
