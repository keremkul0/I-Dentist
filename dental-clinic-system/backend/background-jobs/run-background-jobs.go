package background_jobs

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type TokenService interface {
	DeleteExpiredTokens(ctx context.Context)
}

type PasswordResetTokenService interface {
	DeleteExpiredTokens(ctx context.Context) error
}

func StartCleanExpiredJwtTokens(tokenService TokenService) {
	ctx := context.Background()

	c := cron.New()
	cronExpression := fmt.Sprintf("@every %ds", 10)

	_, err := c.AddFunc(cronExpression, func() {
		tokenService.DeleteExpiredTokens(ctx)
	})
	if err != nil {
		panic(err)
	}

	c.Start()
}

func StartCleanExpiredPasswordResetTokens(passwordResetTokenService PasswordResetTokenService) {
	ctx := context.Background()

	c := cron.New()
	// Her gün gece yarısı çalışsın
	cronExpression := "0 0 * * *"

	_, err := c.AddFunc(cronExpression, func() {
		err := passwordResetTokenService.DeleteExpiredTokens(ctx)
		if err != nil {
			// Log the error but don't panic
			fmt.Printf("Error deleting expired password reset tokens: %v\n", err)
		}
	})
	if err != nil {
		panic(err)
	}

	c.Start()
}
