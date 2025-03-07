package wl

import (
	"cmp"
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"slices"
	"sync"
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

	mu        sync.Mutex
	expansion *PartitionedExpansion
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
func (w *WarcraftLogsClient) GetExpansionEncounters(ctx context.Context) (*PartitionedExpansion, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.expansion != nil {
		return w.expansion, nil
	}

	eeq := &ExpansionEncountersQuery{}
	if err := w.Query(ctx, eeq, nil); err != nil {
		w.l.Error("ExpansionEncountersQuery failed", "error", err)
		return nil, err
	}

	// Sort the expansions and get the latest.
	pe := PartitionedExpansionFromExpansionEncounterQuery(eeq)
	slices.SortFunc(pe, func(a, b *PartitionedExpansion) int {
		return -cmp.Compare(a.ID, b.ID)
	})
	w.expansion = pe[0]

	return w.expansion, nil
}

// GetParsesForCharacter gets the parses for the given character, partition and zone id.
func (w *WarcraftLogsClient) GetParsesForCharacter(ctx context.Context, options *CharacterParsesQueryOptions) (*CharacterParsesQuery, error) {
	var partition int
	if options.Partition == nil {
		// Get the default partition for the zone
		partition = w.GetDefaultPartitionByZoneID(options.ZoneID)
	} else {
		partition = *options.Partition
	}

	cpq := &CharacterParsesQuery{}
	vars := map[string]any{
		"name":          graphql.String(options.Name),
		"server_slug":   graphql.String(options.ServerSlug),
		"server_region": graphql.String(options.ServerRegion),
		"partition":     graphql.Int(partition),
		"zone_id":       graphql.Int(options.ZoneID),
	}
	if err := w.Query(ctx, cpq, vars); err != nil {
		w.l.Error("CharacterParsesQuery failed", "error", err)
		return nil, err
	}
	return cpq, nil
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
		var ne graphql.NetworkError
		if errors.As(err, &ne) && ne.StatusCode() == http.StatusTooManyRequests {
			w.limiter.SpendAllPoints()
			w.l.Warn("WarcraftLogs Rate Limit Exceeded, spent all remaining points")
			return err
		}

		w.l.Error("GraphQL Query errored", "error", err)
		return err
	}

	w.limiter.SetPointsSpent(query.Data())

	return nil
}

// GetDefaultPartitionByZoneID returns the default partition for the give zoneID
func (w *WarcraftLogsClient) GetDefaultPartitionByZoneID(zoneID int) int {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, zone := range w.expansion.Zones {
		if zone.ID == zoneID {

			for _, partition := range zone.Partitions {
				if partition.Default {
					return partition.ID
				}
			}
		}
	}

	return -1
}

// ClearPartitionedExpansion sets the expansions to nil, allowing for it to be re-cached.
func (w *WarcraftLogsClient) ClearPartitionedExpansion() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.expansion = nil
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
		expansion: nil,
	}

	if _, err := wlc.GetRateLimit(context.Background()); err != nil {
		wlc.l.Error("Failed to get RateLimit during setup")
		panic(err)
	}

	if _, err := wlc.GetExpansionEncounters(context.Background()); err != nil {
		wlc.l.Error("Failed to get ExpansionEncounters during setup")
		panic(err)
	}

	return wlc
}
