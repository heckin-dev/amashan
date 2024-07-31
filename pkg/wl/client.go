package wl

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"time"
)

const (
	WL_API_URL = "https://www.warcraftlogs.com/api/v2/client"
)

// WarcraftLogsClient wraps the WarcraftLogs v2 GraphQL API abstracting requests we care about.
type WarcraftLogsClient struct {
	l hclog.Logger

	config  *clientcredentials.Config
	limiter *PointLimiter
}

// GetRateLimit returns the current rate limit remaining for the client.
func (w *WarcraftLogsClient) GetRateLimit(ctx context.Context) (*RateLimitQuery, error) {
	rlq := &RateLimitQuery{}
	if err := w.Query(ctx, rlq, nil); err != nil {
		w.l.Error("RateLimitQuery failed", "error", err)
		return nil, err
	}

	return rlq, nil
}

// GetExpansionEncounters gets the encounters for all the expansions.
func (w *WarcraftLogsClient) GetExpansionEncounters(ctx context.Context) (*ExpansionEncountersQuery, error) {
	eeq := &ExpansionEncountersQuery{}
	if err := w.Query(ctx, eeq, nil); err != nil {
		w.l.Error("ExpansionEncountersQuery failed", "error", err)
		return nil, err
	}

	return eeq, nil
}

func (w *WarcraftLogsClient) GetPartitionsForExpansion(ctx context.Context, expansionID int) (*ExpansionPartitionsQuery, error) {
	epq := &ExpansionPartitionsQuery{}
	vars := map[string]any{
		"expansion_id": graphql.Int(expansionID),
	}
	if err := w.Query(ctx, epq, vars); err != nil {
		w.l.Error("ExpansionPartitionsQuery failed", "error", err)
		return nil, err
	}
	return epq, nil
}

// Query performs a query.
func (w *WarcraftLogsClient) Query(ctx context.Context, query RatedQuery, vars map[string]interface{}) error {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
	}

	if err := w.limiter.CanSpendPoints(); err != nil {
		return err
	}

	client := graphql.NewClient(WL_API_URL, w.config.Client(ctx))
	if err := client.Query(ctx, query, vars); err != nil {
		if errors.Is(err, graphql.ErrClientResponse) {
			w.limiter.SpendAllPoints()
			w.l.Warn("WarcraftLogs Query received a 4xx Status Code")
			return err
		}

		w.l.Error("GraphQL Query errored", "error", err)
		return err
	}

	w.limiter.SetPointsSpent(query.Data())

	return nil
}

// NewWarcraftLogsClient creates a new client and gets the remaining rate-limit.
func NewWarcraftLogsClient(l hclog.Logger) *WarcraftLogsClient {
	wlc := &WarcraftLogsClient{
		l: l,
		config: &clientcredentials.Config{
			ClientID:     os.Getenv("WL_CLIENT_ID"),
			ClientSecret: os.Getenv("WL_CLIENT_SECRET"),
			TokenURL:     "https://www.warcraftlogs.com/oauth/token",
		},
		limiter: NewPointLimiter(l, RateLimitData{
			LimitPerHour:        3600,
			PointsSpentThisHour: 0,
			PointsResetIn:       3600,
		}),
	}

	if _, err := wlc.GetRateLimit(context.Background()); err != nil {
		wlc.l.Error("Failed to get RateLimit during setup")
		panic(err)
	}

	return wlc
}
