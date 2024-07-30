package wl

import (
	"github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
	"os"
	"time"
)

type WarcraftLogsClient struct {
	l hclog.Logger

	config         *clientcredentials.Config
	perHourLimiter *rate.Limiter
}

/*
	TODO: We need to create our client.
	TODO: We need to make a request to the endpoint.
	TODO: We need to check for any errors.
	TODO: We need to check that we aren't exceeding the Token Limit.
	TODO: We should assume if we have Tokens <= 5, we just don't make the request.

	! This must be thread-safe.

	// 1. Make a request to get the Rate Limit
	// 2. Store the Rate Limit on the Client (limitPerHour, pointsResetIn)
	// 3. Start a 1-shot timer for the Reset limit.
	// 4. Store pointsSpentThisHour
	// 5. If we get an error about RateLimit, block all requests until future.
	// 6. If we have no points to spend, don't do anything.
*/

func NewWarcraftLogsClient(l hclog.Logger) *WarcraftLogsClient {
	return &WarcraftLogsClient{
		l: l,
		config: &clientcredentials.Config{
			ClientID:     os.Getenv("WL_CLIENT_ID"),
			ClientSecret: os.Getenv("WL_CLIENT_SECRET"),
			TokenURL:     "https://www.warcraftlogs.com/oauth/token",
		},
		perHourLimiter: rate.NewLimiter(rate.Every(1*time.Hour), 3600),
	}
}
