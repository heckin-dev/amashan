package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/bnet"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"github.com/heckin-dev/amashan/pkg/utils"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

type BattleNet struct {
	l hclog.Logger

	client *bnet.BattlenetClient
	store  *sessions.CookieStore
}

func (b *BattleNet) Authorize(w http.ResponseWriter, r *http.Request) {
	state, err := utils.NewStateString(64)
	if err != nil {
		b.l.Error("failed to generate state string", "error", err)
		http.Error(w, "failed to generate state string", http.StatusInternalServerError)
		return
	}

	// Create a new Session and store the state.
	session, _ := b.store.Get(r, "oauth")
	session.Values["state"] = state

	// Save the session
	if err := session.Save(r, w); err != nil {
		b.l.Error("failed to save session", "error", err)
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, b.client.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (b *BattleNet) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := b.store.Get(r, "oauth")
	if err != nil || session == nil || session.IsNew {
		http.Error(w, "no session found for this request", http.StatusBadRequest)
		return
	}

	// Get the initial request state
	rState, ok := session.Values["state"].(string)
	if !ok {
		http.Error(w, "failed to read state from session", http.StatusBadRequest)
		return
	}

	// Ensure the callback state is equal to the request state
	cbState := r.URL.Query().Get("state")
	if !strings.EqualFold(rState, cbState) {
		http.Error(w, "callback state mismatch", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := b.client.Exchange(r.Context(), code)
	if err != nil {
		b.l.Error("token exchange failed", "error", err)
		http.Error(w, "failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// Check the token.
	ct, err := b.client.CheckToken(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ui, err := b.client.UserInfo(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b.l.Info("callback", "check_token", ct, "userinfo", ui)

	// TODO: Something with the token?
	// Access token
	// Token Type
	// Expiry

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

func (b *BattleNet) ProfileSummary(w http.ResponseWriter, r *http.Request) {
	region := r.Context().Value(middleware.RegionContextKey).(string)

	// TODO: Read the Token from the request header or something.
	// TODO: Add the token to the line below.
	as, err := b.client.AccountProfileSummary(r.Context(), &oauth2.Token{}, region)
	if err != nil {
		b.l.Error("failed to retrieve account summary", "error", err)
		http.Error(w, "failed to retrieve account summary", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(as)
}

func (b *BattleNet) Route(r *mux.Router) {
	oauthRouter := r.PathPrefix("/auth").Subrouter()

	// http://localhost:9090/api/auth/battlenet
	oauthRouter.HandleFunc("/battlenet", b.Authorize).Methods(http.MethodGet)
	oauthRouter.HandleFunc("/battlenet/callback", b.Callback).Methods(http.MethodGet)

	bnetRouter := r.PathPrefix("/{region}/battlenet").Subrouter()
	bnetRouter.Use(middleware.UseRegion().Middleware)

	// http://localhost:9090/api/us/battlenet/profile
	bnetRouter.HandleFunc("/profile", b.ProfileSummary).Methods(http.MethodGet)
}

func NewBattleNet(l hclog.Logger) *BattleNet {
	// Create a new store
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = !strings.EqualFold(os.Getenv("PROD"), "")

	return &BattleNet{
		l:      l,
		client: bnet.NewBattlnetClient(l),
		store:  store,
	}
}
