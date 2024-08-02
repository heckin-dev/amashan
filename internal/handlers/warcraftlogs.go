package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"github.com/heckin-dev/amashan/pkg/wl"
	"net/http"
	"strconv"
	"time"
)

type WarcraftLogs struct {
	l hclog.Logger

	client *wl.WarcraftLogsClient
}

func (wls *WarcraftLogs) ClearCachedExpansion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	if len(r.Header.Values("X-Amashan-Anonymous-Authority")) == 0 {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// TODO: Might be worth changing this magic value to something we know always.
	cache := r.Context().Value(middleware.CacheContextKey).(middleware.CacheClient)
	go cache.Del("/api/warcraftlogs/partitions")

	wls.client.ClearPartitionedExpansion()
	w.WriteHeader(http.StatusNoContent)
}

func (wls *WarcraftLogs) Partitions(w http.ResponseWriter, r *http.Request) {
	cache := r.Context().Value(middleware.CacheContextKey).(middleware.CacheClient)
	key := r.URL.Path

	// Cache HIT
	if val, err := cache.Get(r.Context(), key); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(val))
		return
	}

	partitionedExpansion, err := wls.client.GetExpansionEncounters(r.Context())
	if err != nil {
		wls.l.Error("failed to retrieve partitioned expansion", "error", err)
		http.Error(w, "failed to retrieve partitioned expansion", http.StatusInternalServerError)
		return
	}

	// Marshal the PartitionedExpansion
	bs, err := json.Marshal(partitionedExpansion)
	if err != nil {
		wls.l.Error("json.Marshal failed for PartitionedExpansion", "error", err)
		http.Error(w, "failed to marshal partitioned expansion", http.StatusInternalServerError)
		return
	}

	// Cache SET
	go cache.Set(key, string(bs), 12*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
}

func (wls *WarcraftLogs) CharacterParses(w http.ResponseWriter, r *http.Request) {
	cache := r.Context().Value(middleware.CacheContextKey).(middleware.CacheClient)
	key := r.RequestURI

	// Cache HIT
	if val, err := cache.Get(r.Context(), key); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(val))
		return
	}

	options := wl.CharacterParsesQueryOptionsFromContext(r.Context())

	q := r.URL.Query()
	if !q.Has("zone_id") {
		http.Error(w, "missing required query param 'zone_id'", http.StatusBadRequest)
		return
	}

	zoneStr := q.Get("zone_id")
	zone, err := strconv.Atoi(zoneStr)
	if err != nil {
		http.Error(w, "query param 'zone_id' must be an integer", http.StatusBadRequest)
		return
	}
	options.ZoneID = zone

	if q.Has("partition") {
		partitionStr := q.Get("partition")
		partition, err := strconv.Atoi(partitionStr)
		if err != nil {
			http.Error(w, "optional query param 'partition' must be an integer", http.StatusBadRequest)
			return
		}
		options.Partition = &partition
	}

	parses, err := wls.client.GetParsesForCharacter(r.Context(), options)
	if err != nil {
		wls.l.Error("failed to retrieve character parses", "error", err)
		http.Error(w, "failed to retrieve character parses", http.StatusInternalServerError)
		return
	}

	// Marshal the CharacterParses
	bs, err := json.Marshal(parses)
	if err != nil {
		wls.l.Error("json.Marshal failed for CharacterParses", "error", err)
		http.Error(w, "failed to marshal character parses", http.StatusInternalServerError)
		return
	}

	// Cache SET
	go cache.Set(key, string(bs), 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bs)
}

func (wls *WarcraftLogs) Route(r *mux.Router) {
	wlRouter := r.PathPrefix("/warcraftlogs").Subrouter()

	wlRouter.HandleFunc("", wls.ClearCachedExpansion)
	wlRouter.HandleFunc("/partitions", wls.Partitions)

	rrcRouter := wlRouter.PathPrefix("/{region}/{realm}/{character}").Subrouter()
	rrcRouter.Use(middleware.UseRegion().Middleware)
	rrcRouter.Use(middleware.UseRealm().Middleware)
	rrcRouter.Use(middleware.UseCharacter().Middleware)

	rrcRouter.HandleFunc("/parses", wls.CharacterParses)
}

func NewWarcraftLogs(l hclog.Logger) *WarcraftLogs {
	return &WarcraftLogs{
		l:      l,
		client: wl.NewWarcraftLogsClient(l),
	}
}
