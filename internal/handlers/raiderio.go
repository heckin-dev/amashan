package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"github.com/heckin-dev/amashan/pkg/rio"
	"net/http"
	"time"
)

type RaiderIO struct {
	l      hclog.Logger
	client *rio.RaiderIOClient
}

func (i *RaiderIO) CharacterProfile(w http.ResponseWriter, r *http.Request) {
	cache := r.Context().Value(middleware.CacheContextKey).(middleware.CacheClient)
	key := r.URL.Path

	// Cache HIT
	if val, err := cache.Get(r.Context(), key); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(val))
		return
	}

	profile, err := i.client.CharacterProfile(r.Context(), rio.CharacterProfileOptionsFromContext(r.Context()))
	if err != nil {
		i.l.Error("failed to retrieve raiderio character profile", "error", err)
		http.Error(w, "failed to retrieve raiderio character profile", http.StatusInternalServerError)
		return
	}

	// Marshal the CharacterProfile
	bs, err := json.Marshal(profile)
	if err != nil {
		i.l.Error("json.Marshal failed for CharacterProfile", "error", err)
		http.Error(w, "failed to marshal raiderio character profile", http.StatusInternalServerError)
		return
	}

	// Cache SET
	go cache.Set(key, string(bs), 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
}

func (i *RaiderIO) Route(r *mux.Router) {
	rioRouter := r.PathPrefix("/raiderio/{region}/{realm}/{character}").Subrouter()
	rioRouter.Use(middleware.UseRegion().Middleware)
	rioRouter.Use(middleware.UseRealm().Middleware)
	rioRouter.Use(middleware.UseCharacter().Middleware)

	rioRouter.HandleFunc("", i.CharacterProfile)
}

func NewRaiderIO(l hclog.Logger) *RaiderIO {
	return &RaiderIO{
		l:      l,
		client: rio.NewRaiderIOClient(l),
	}
}
