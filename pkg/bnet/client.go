package bnet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	BNET_OAUTH_URL string = "https://oauth.battle.net"
	BNET_API_URL          = "https://{region}.api.blizzard.com"
)

// RegionalURLFunc wraps the string replacement for building the Region-ed API URL making it more testable.
type RegionalURLFunc func(region string) string

type BattlenetClient struct {
	l hclog.Logger

	clientConfig     *cc.Config
	oauthConfig      *oauth2.Config
	perSecondLimiter *rate.Limiter
	perHourLimiter   *rate.Limiter

	apiURLFn RegionalURLFunc
}

// AuthCodeURL returns the AuthCodeURL produced by the underlying oauth2.Config to be redirected to for OAuth2.
func (b *BattlenetClient) AuthCodeURL(state string) string {
	return b.oauthConfig.AuthCodeURL(state)
}

// Exchange returns the *oauth2.Token and/or error produced during the token exchange.
func (b *BattlenetClient) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return b.oauthConfig.Exchange(ctx, code)
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
	res, err := b.Do(ctx, t, req, OAuthRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
	res, err := b.Do(ctx, t, req, OAuthRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode the response.
	uiRes := &UserInfoResponse{}
	err = json.NewDecoder(res.Body).Decode(uiRes)
	if err != nil {
		return nil, err
	}

	return uiRes, nil
}

// AccountProfileSummary gets the account summary for the given token.
func (b *BattlenetClient) AccountProfileSummary(ctx context.Context, options *AccountSummaryOptions) (*AccountSummaryResponse, error) {
	const endpoint string = "/profile/user/wow"

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	options.Token.SetAuthHeader(req)

	res, err := b.Do(ctx, options.Token, req, OAuthRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	asRes := &AccountSummaryResponse{}
	err = json.NewDecoder(res.Body).Decode(asRes)
	if err != nil {
		return nil, err
	}

	return asRes, nil
}

// CharacterSummary gets the summary for a given character.
func (b *BattlenetClient) CharacterSummary(ctx context.Context, options *CharacterOptions) (*CharacterSummaryResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	csRes := &CharacterSummaryResponse{}
	err = json.NewDecoder(res.Body).Decode(csRes)
	if err != nil {
		return nil, err
	}

	return csRes, nil
}

// CharacterStatus gets the status for a given character.
func (b *BattlenetClient) CharacterStatus(ctx context.Context, options *CharacterOptions) (*CharacterStatusResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/status", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	csRes := &CharacterStatusResponse{}
	err = json.NewDecoder(res.Body).Decode(csRes)
	if err != nil {
		return nil, err
	}

	return csRes, nil
}

// CharacterEquipmentSummary gets the equipment summary for a given character.
func (b *BattlenetClient) CharacterEquipmentSummary(ctx context.Context, options *CharacterOptions) (*CharacterEquipmentResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/equipment", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ceRes := &CharacterEquipmentResponse{}
	err = json.NewDecoder(res.Body).Decode(ceRes)
	if err != nil {
		return nil, err
	}

	return ceRes, nil
}

// CharacterMedia gets the character media for a given character.
func (b *BattlenetClient) CharacterMedia(ctx context.Context, options *CharacterOptions) (*CharacterMediaResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/character-media", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cmRes := &CharacterMediaResponse{}
	err = json.NewDecoder(res.Body).Decode(cmRes)
	if err != nil {
		return nil, err
	}

	return cmRes, nil
}

// CharacterStatistics gets the character statistics for a given character.
func (b *BattlenetClient) CharacterStatistics(ctx context.Context, options *CharacterOptions) (*CharacterStatisticsResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/statistics", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	csRes := &CharacterStatisticsResponse{}
	err = json.NewDecoder(res.Body).Decode(csRes)
	if err != nil {
		return nil, err
	}

	return csRes, nil
}

// CharacterDungeonEncounters gets the dungeon encounters for the given character.
func (b *BattlenetClient) CharacterDungeonEncounters(ctx context.Context, options *CharacterOptions) (*CharacterDungeonEncountersResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/encounters/dungeons", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cdeRes := &CharacterDungeonEncountersResponse{}
	err = json.NewDecoder(res.Body).Decode(cdeRes)
	if err != nil {
		return nil, err
	}

	return cdeRes, nil
}

// CharacterRaidEncounters gets the raid encounters for the given character.
func (b *BattlenetClient) CharacterRaidEncounters(ctx context.Context, options *CharacterOptions) (*CharacterRaidEncountersResponse, error) {
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/encounters/raids", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	creRes := &CharacterRaidEncountersResponse{}
	err = json.NewDecoder(res.Body).Decode(creRes)
	if err != nil {
		return nil, err
	}

	return creRes, nil
}

// MythicKeystoneIndex gets the mythic keystone index for the given character.
func (b *BattlenetClient) MythicKeystoneIndex(ctx context.Context, options *CharacterOptions) (*MythicKeystoneIndexResponse, error) {
	// /profile/wow/character/{realmSlug}/{characterName}/mythic-keystone-profile
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/mythic-keystone-profile", options.Realm, options.Character)

	req, err := b.prepareRequest(&RequestOptions{
		Region:    options.Region,
		Namespace: ProfileNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	mkiRes := &MythicKeystoneIndexResponse{}
	err = json.NewDecoder(res.Body).Decode(mkiRes)
	if err != nil {
		return nil, err
	}

	return mkiRes, nil
}

// MythicKeystoneSeason gets the mythic keystone season for the given character.
func (b *BattlenetClient) MythicKeystoneSeason(ctx context.Context, options *MythicSeasonOptions) (*MythicKeystoneSeasonResponse, error) {
	// /profile/wow/character/{realmSlug}/{characterName}/mythic-keystone-profile/season/{seasonId}
	var endpoint = fmt.Sprintf("/profile/wow/character/%s/%s/mythic-keystone-profile/season/%d", options.Realm, options.Character, options.Season)

	req, err := b.prepareRequest(&RequestOptions{
		Region:      options.Region,
		Namespace:   ProfileNamespace,
		Endpoint:    endpoint,
		Method:      http.MethodGet,
		QueryParams: map[string]string{"seasonId": strconv.Itoa(options.Season)},
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		// If we got a 404, this is expected, we should handle with value.
		var errUnexpectedResponse *ErrUnexpectedResponse
		if errors.As(err, &errUnexpectedResponse) {
			return &MythicKeystoneSeasonResponse{
				CharacterPlayedSeason: false,
			}, nil
		}

		return nil, err
	}
	defer res.Body.Close()

	mksRes := &MythicKeystoneSeasonResponse{}
	err = json.NewDecoder(res.Body).Decode(mksRes)
	if err != nil {
		return nil, err
	}

	mksRes.CharacterPlayedSeason = true
	return mksRes, nil
}

// RealmsByRegion gets the realm index for the given region.
func (b *BattlenetClient) RealmsByRegion(ctx context.Context, region RegionOption) (*RealmIndexResponse, error) {
	// /data/wow/realm/index
	const endpoint = "/data/wow/realm/index"

	req, err := b.prepareRequest(&RequestOptions{
		Region:    region.String(),
		Namespace: DynamicNamespace,
		Endpoint:  endpoint,
		Method:    http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	res, err := b.Do(ctx, nil, req, ClientRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	riRes := &RealmIndexResponse{}
	err = json.NewDecoder(res.Body).Decode(riRes)
	if err != nil {
		return nil, err
	}

	riRes.Region = region.String()
	return riRes, nil
}

// Do does the provided *http.Request using the http.Client associated with the provided *oauth2.Token. This can be
// used directly but there are likely other wrapper methods that are more useful.
func (b *BattlenetClient) Do(ctx context.Context, t *oauth2.Token, req *http.Request, rType RequestType) (*http.Response, error) {
	// If we are making a request on behalf of the user, we should check that their token is valid.
	if rType == OAuthRequest {
		if !t.Valid() {
			return nil, ErrTokenIsInvalid
		}
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

	var res *http.Response
	var err error

	if rType == ClientRequest {
		res, err = b.clientConfig.Client(ctx).Do(req)
	} else {
		res, err = b.oauthConfig.Client(ctx, t).Do(req)
	}

	if err != nil {
		b.l.Error("Failed to do request", "RequestType", rType, "request", req, "error", err)
		return nil, err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		b.perSecondLimiter.ReserveN(time.Now().Add(1*time.Minute), b.perSecondLimiter.Burst())
		b.perHourLimiter.ReserveN(time.Now().Add(1*time.Hour), b.perHourLimiter.Burst())
		b.l.Info("BattleNet Rate-Limit reached, drained remaining tokens")
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			b.l.Error("Failed to read body for error logging", "error", err)
		}
		defer res.Body.Close()

		if res.StatusCode == 404 {
			b.l.Warn("Request returned 404 status code", "body", string(bs))
		} else {
			b.l.Error("Request returned non-200 status code", "StatusCode", res.StatusCode, "body", string(bs))
		}

		return nil, &ErrUnexpectedResponse{StatusCode: res.StatusCode, Err: err}
	}

	return res, nil
}

// SetConfig overrides the underlying oauth2.Config with the provided one.
//
//	Should only be used for testing.
//	Let the environment variables configure it otherwise.
func (b *BattlenetClient) SetConfig(config *oauth2.Config) {
	b.oauthConfig = config
}

// prepareRequest util wraps common http.NewRequest(...) and query param setup.
//
// by default this provides the following query params:
//
//	?region=RequestOptions.Region
//	&namespaced=RequestOptions.Namespace-RequestOptions.Region
//	&locale=en_US
func (b *BattlenetClient) prepareRequest(options *RequestOptions) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", b.apiURLFn(options.Region), options.Endpoint)
	req, err := http.NewRequest(options.Method, url, options.Body)
	if err != nil {
		b.l.Error("failed to create request", "url", url, "error", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region", options.Region)
	q.Add("namespace", fmt.Sprintf("%s-%s", options.Namespace, options.Region))
	q.Add("locale", "en_US")

	if options.QueryParams != nil {
		for k, v := range options.QueryParams {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func NewBattlnetClient(l hclog.Logger) *BattlenetClient {
	return &BattlenetClient{
		l: l,
		clientConfig: &cc.Config{
			ClientID:     os.Getenv("BNET_CLIENT_ID"),
			ClientSecret: os.Getenv("BNET_CLIENT_SECRET"),
			TokenURL:     "https://oauth.battle.net/token",
		},
		oauthConfig: &oauth2.Config{
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
		apiURLFn: func(region string) string {
			return strings.Replace(BNET_API_URL, "{region}", region, -1)
		},
	}
}
