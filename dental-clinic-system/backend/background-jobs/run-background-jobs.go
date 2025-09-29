package background_jobs

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type TokenService interface {
	DeleteExpiredTokens(ctx context.Context)
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
