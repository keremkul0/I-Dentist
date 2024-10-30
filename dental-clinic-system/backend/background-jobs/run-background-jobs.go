package background_jobs

import (
	"dental-clinic-system/application/tokenService"
	"fmt"
	"github.com/robfig/cron/v3"
)

func StartCleanExpiredJwtTokens(tokenService tokenService.TokenService) {
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
