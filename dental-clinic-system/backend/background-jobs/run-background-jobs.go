package background_jobs

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type TokenService interface {
	DeleteExpiredTokensService()
}

func StartCleanExpiredJwtTokens(tokenService TokenService) {
	c := cron.New()
	cronExpression := fmt.Sprintf("@every %ds", 10)

	_, err := c.AddFunc(cronExpression, func() {
		tokenService.DeleteExpiredTokensService()
	})
	if err != nil {
		panic(err)
	}

	c.Start()
}
