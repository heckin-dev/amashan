package bnet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

/*
	TODO: Create a Battlnet API Wrapper Client

	https://develop.battle.net/documentation/guides/getting-started

	1. The client should include a rate-limiter / round-tripper
	2. We should make it easy to strap up a token for a given request.
	3. Providing a middleware might be useful for strapping up the client between requests.

	: Flow?

	1. User does login
	2. * We retrieve info about the user identity? (do we care)
	3. Get the user's characters
	4. ! Determine how we retrieve the appearances the user has?
	5. Get the user's auction house data
	6. Show the user the appearances on the AH they don't have


	--

	Determine the items and their appearances?
	Determine the items that share their appearances.
	Determine the items that can be purchased on the auction house.
	...


https://stackoverflow.com/questions/51628755/how-to-add-default-header-fields-from-http-client
https://medium.com/mflow/rate-limiting-in-golang-http-client-a22fba15861a
*/

const (
	BNET_OAUTH_URL string = "https://oauth.battle.net"
	BNET_API_URL          = "https://{region}.api.blizzard.com"
)

type BattlenetClient struct {
	l hclog.Logger

	config           *oauth2.Config
	perSecondLimiter *rate.Limiter
	perHourLimiter   *rate.Limiter
}

// AuthCodeURL returns the AuthCodeURL produced by the underlying oauth2.Config to be redirected to for OAuth2.
func (b *BattlenetClient) AuthCodeURL(state string) string {
	return b.config.AuthCodeURL(state)
}

// Exchange returns the *oauth2.Token and/or error produced during the token exchange.
func (b *BattlenetClient) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return b.config.Exchange(ctx, code)
}

// CheckToken ensures the token is valid and contains the required scopes.
func (b *BattlenetClient) CheckToken(ctx context.Context, t *oauth2.Token) (*CheckTokenResponse, error) {
	const endpoint string = "/oauth/check_token"

	// Create the request.
	url := fmt.Sprintf("%s%s", BNET_OAUTH_URL, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		b.l.Error("failed to create request", "url", url, "error", err)
		return nil, err
	}

	// Add the required query params.
	q := req.URL.Query()

	q.Add("region", "us")
	q.Add("token", t.AccessToken)

	req.URL.RawQuery = q.Encode()

	// Do the request.
	res, err := b.Do(ctx, t, req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	// Decode the response and validate.
	ctRes := &CheckTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(ctRes)
	if err != nil {
		return nil, err
	}

	// Ensure we have the required scopes.
	if !slices.Contains(ctRes.Scope, "wow.profile") {
		return nil, ErrMissingRequiredScope{Scope: "wow.profile"}
	}

	if !slices.Contains(ctRes.Scope, "openid") {
		return nil, ErrMissingRequiredScope{Scope: "openid"}
	}

	return ctRes, nil
}

// UserInfo gets the userinfo for the given token.
func (b *BattlenetClient) UserInfo(ctx context.Context, t *oauth2.Token) (*UserInfoResponse, error) {
	const endpoint string = "/oauth/userinfo"

	// Create the request.
	url := fmt.Sprintf("%s%s", BNET_OAUTH_URL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		b.l.Error("failed to create request", "url", url, "error", err)
		return nil, err
	}

	// Add the required query params.
	q := req.URL.Query()

	q.Add("region", "us")

	req.URL.RawQuery = q.Encode()

	// Add the token to the header.
	t.SetAuthHeader(req)

	// Do the request.
	res, err := b.Do(ctx, t, req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	// Decode the response.
	uiRes := &UserInfoResponse{}
	err = json.NewDecoder(res.Body).Decode(uiRes)
	if err != nil {
		return nil, err
	}

	return uiRes, nil
}

// AccountProfileSummary gets the account summary for the given token.
func (b *BattlenetClient) AccountProfileSummary(ctx context.Context, t *oauth2.Token, region string) (*AccountSummaryResponse, error) {
	const endpoint string = "/profile/user/wow"

	url := fmt.Sprintf("%s%s", strings.Replace(BNET_API_URL, "{region}", region, -1), endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		b.l.Error("failed to create request", "url", url, "error", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", region)
	q.Add("namespace", fmt.Sprintf("profile-%s", region))
	q.Add("locale", "en_US")
	req.URL.RawQuery = q.Encode()

	t.SetAuthHeader(req)

	res, err := b.Do(ctx, t, req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	asRes := &AccountSummaryResponse{}
	err = json.NewDecoder(res.Body).Decode(asRes)
	if err != nil {
		return nil, err
	}

	return asRes, nil
}

// Do does the provided *http.Request using the http.Client associated with the provided *oauth2.Token. This can be
// used directly but there are likely other wrapper methods that are more useful.
func (b *BattlenetClient) Do(ctx context.Context, t *oauth2.Token, req *http.Request) (*http.Response, error) {
	if !t.Valid() {
		return nil, ErrTokenIsInvalid
	}

	// If we create a new context, we need to defer it.
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
	}

	// Ensure we aren't exceeding the hourly rate limit.
	if err := b.perHourLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	// Ensure we aren't exceeding the per-second rate limit.
	if err := b.perSecondLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return b.config.Client(ctx, t).Do(req)
}

// SetConfig overrides the underlying oauth2.Config with the provided one.
//
//	Should only be used for testing.
//	Let the environment variables configure it otherwise.
func (b *BattlenetClient) SetConfig(config *oauth2.Config) {
	b.config = config
}

func NewBattlnetClient(l hclog.Logger) *BattlenetClient {
	return &BattlenetClient{
		l: l,
		config: &oauth2.Config{
			ClientID:     os.Getenv("BNET_CLIENT_ID"),
			ClientSecret: os.Getenv("BNET_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.battle.net/authorize",
				TokenURL: "https://oauth.battle.net/token",
			},
			RedirectURL: os.Getenv("BNET_REDIRECT_URL"),
			Scopes:      []string{"wow.profile", "openid"},
		},
		perSecondLimiter: rate.NewLimiter(rate.Every(1*time.Second), 100), // 100/s
		perHourLimiter:   rate.NewLimiter(rate.Every(1*time.Hour), 36000), // 36,000/h
	}
}
