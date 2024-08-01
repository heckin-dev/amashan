package rio

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"time"
)

const (
	API_URL = "https://raider.io/api/v1"
)

// URLFunc wraps the string for the api url to allow for testing.
type URLFunc func() string

// RaiderIOClient wraps the RaiderIO API
type RaiderIOClient struct {
	l                hclog.Logger
	perMinuteLimiter *rate.Limiter
	apiURLFn         URLFunc
}

// CharacterProfile gets a character's mythic plus statistics.
func (r *RaiderIOClient) CharacterProfile(ctx context.Context, options *CharacterProfileOptions) (*CharacterProfileResponse, error) {
	const endpoint = "/characters/profile"
	const fields = "mythic_plus_ranks,mythic_plus_recent_runs,mythic_plus_best_runs,mythic_plus_scores_by_season:current"

	url := fmt.Sprintf("%s%s", r.apiURLFn(), endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		r.l.Error("Failed to create request", "url", url, "error", err)
		return nil, err
	}

	// Set query params
	q := req.URL.Query()
	q.Add("region", options.Region)
	q.Add("realm", options.Realm)
	q.Add("name", options.Character)
	q.Add("fields", fields)
	req.URL.RawQuery = q.Encode()

	res, err := r.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cpRes := &CharacterProfileResponse{}
	err = json.NewDecoder(res.Body).Decode(cpRes)
	if err != nil {
		return nil, err
	}

	return cpRes, nil
}

// Do handles making http requests ensuring they abide by the given rate-limits.
func (r *RaiderIOClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// If we create a new context, we need to defer it.
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
	}

	// Ensure we aren't exceeding the minute rate limit.
	if err := r.perMinuteLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	// Make the request
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		r.l.Error("Failed to do request", "request", req, "error", err)
		return nil, err
	}

	// If we got a 404, we should drain the remaining tokens.
	if res.StatusCode == http.StatusTooManyRequests {
		r.perMinuteLimiter.ReserveN(time.Now().Add(1*time.Minute), r.perMinuteLimiter.Burst())
		r.l.Info("RaiderIO Rate-Limit reached, drained remaining tokens")
	}

	// If we didn't get an OK, we should just error.
	if res.StatusCode != http.StatusOK {
		bs, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			r.l.Error("Failed to read body for error logging", "error", err)
			return nil, err
		}

		r.l.Error("RaiderIO Response was non-200", "error", err, "body", string(bs))
		return nil, fmt.Errorf("non-200 status code: '%d', body: '%s'", res.StatusCode, string(bs))
	}

	return res, nil
}

// NewRaiderIOClient creates a new default RaiderIOClient
func NewRaiderIOClient(l hclog.Logger) *RaiderIOClient {
	return &RaiderIOClient{
		l:                l,
		perMinuteLimiter: rate.NewLimiter(rate.Every(1*time.Minute), 300),
		apiURLFn: func() string {
			return API_URL
		},
	}
}
