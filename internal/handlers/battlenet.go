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
	"strconv"
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
	session, err := b.store.Get(r, "oauth")
	if err != nil {
		b.l.Error("failed to decode existing session", "error", err)
		http.Error(w, "failed to decode existing session", http.StatusInternalServerError)
		return
	}
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
	// TODO: Read the Token from the request header or something.
	// TODO: Add the token to the line below.
	as, err := b.client.AccountProfileSummary(r.Context(), &bnet.AccountSummaryOptions{
		Token:  &oauth2.Token{},
		Region: r.Context().Value(middleware.RegionContextKey).(string),
	})
	if err != nil {
		b.l.Error("failed to retrieve account summary", "error", err)
		http.Error(w, "failed to retrieve account summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(as)
}

func (b *BattleNet) CharacterSummary(w http.ResponseWriter, r *http.Request) {
	cs, err := b.client.CharacterSummary(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character summary", "error", err)
		http.Error(w, "failed to retrieve character summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cs)
}

func (b *BattleNet) CharacterEquipment(w http.ResponseWriter, r *http.Request) {
	ce, err := b.client.CharacterEquipmentSummary(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character equipment", "error", err)
		http.Error(w, "failed to retrieve character equipment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ce)
}

func (b *BattleNet) CharacterMedia(w http.ResponseWriter, r *http.Request) {
	cm, err := b.client.CharacterMedia(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character media", "error", err)
		http.Error(w, "failed to retrieve character media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cm)
}

func (b *BattleNet) CharacterStatistics(w http.ResponseWriter, r *http.Request) {
	cs, err := b.client.CharacterStatistics(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character statistics", "error", err)
		http.Error(w, "failed to retrieve character statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cs)
}

func (b *BattleNet) CharacterDungeonEncounters(w http.ResponseWriter, r *http.Request) {
	cde, err := b.client.CharacterDungeonEncounters(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character dungeon encounters", "error", err)
		http.Error(w, "failed to retrieve character dungeon encounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cde)
}

func (b *BattleNet) CharacterRaidEncounters(w http.ResponseWriter, r *http.Request) {
	cre, err := b.client.CharacterRaidEncounters(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve character raid encounters", "error", err)
		http.Error(w, "failed to retrieve character raid encounters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cre)
}

func (b *BattleNet) MythicKeystoneIndex(w http.ResponseWriter, r *http.Request) {
	mki, err := b.client.MythicKeystoneIndex(r.Context(), bnet.CharacterOptionsFromContext(r.Context()))
	if err != nil {
		b.l.Error("failed to retrieve mythic keystone index", "error", err)
		http.Error(w, "failed to retrieve mythic keystone index", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mki)
}

func (b *BattleNet) MythicKeystoneSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonStr, ok := vars["seasonID"]
	if !ok {
		http.Error(w, "seasonID not provided in route parameter", http.StatusBadRequest)
		return
	}

	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		http.Error(w, "failed to parse seasonID to integer", http.StatusBadRequest)
		return
	}

	mks, err := b.client.MythicKeystoneSeason(r.Context(), &bnet.MythicSeasonOptions{
		CharacterOptions: *bnet.CharacterOptionsFromContext(r.Context()),
		Season:           season,
	})
	if err != nil {
		b.l.Error("failed to retrieve mythic keystone season", "error", err)
		http.Error(w, "failed to retrieve mythic keystone season", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mks)
}

func (b *BattleNet) Route(r *mux.Router) {
	oauthRouter := r.PathPrefix("/auth").Subrouter()

	oauthRouter.HandleFunc("/battlenet", b.Authorize).Methods(http.MethodGet)
	oauthRouter.HandleFunc("/battlenet/callback", b.Callback).Methods(http.MethodGet)

	regionalWowRouter := r.PathPrefix("/{region}/wow").Subrouter()
	regionalWowRouter.Use(middleware.UseRegion().Middleware)

	realmAndCharacterRouter := regionalWowRouter.PathPrefix("/{realm}/{character}").Subrouter()
	realmAndCharacterRouter.Use(middleware.UseRealm().Middleware)
	realmAndCharacterRouter.Use(middleware.UseCharacter().Middleware)

	realmAndCharacterRouter.HandleFunc("", b.CharacterSummary)
	realmAndCharacterRouter.HandleFunc("/equipment", b.CharacterEquipment)
	realmAndCharacterRouter.HandleFunc("/character-media", b.CharacterMedia)
	realmAndCharacterRouter.HandleFunc("/character-statistics", b.CharacterStatistics)
	realmAndCharacterRouter.HandleFunc("/mythic-keystone-index", b.MythicKeystoneIndex)
	realmAndCharacterRouter.HandleFunc("/mythic-keystone-index/season/{seasonID}", b.MythicKeystoneSeason)
	realmAndCharacterRouter.HandleFunc("/encounters/dungeons", b.CharacterDungeonEncounters)
	realmAndCharacterRouter.HandleFunc("/encounters/raids", b.CharacterRaidEncounters)

	/*
		// http://localhost:9090/api/us/wow/profile
		bnetRouter.HandleFunc("/profile", b.ProfileSummary).Methods(http.MethodGet)
	*/
}

func NewBattleNet(l hclog.Logger) *BattleNet {
	// Create a new store
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteNoneMode

	return &BattleNet{
		l:      l,
		client: bnet.NewBattlnetClient(l),
		store:  store,
	}
}
