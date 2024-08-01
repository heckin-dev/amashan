package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"github.com/heckin-dev/amashan/pkg/wl"
	"net/http"
	"strconv"
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

	wls.client.ClearPartitionedExpansion()
	w.WriteHeader(http.StatusNoContent)
}

func (wls *WarcraftLogs) Partitions(w http.ResponseWriter, r *http.Request) {
	partitionedExpansion, err := wls.client.GetExpansionEncounters(r.Context())
	if err != nil {
		wls.l.Error("failed to retrieve partitioned expansion", "error", err)
		http.Error(w, "failed to retrieve partitioned expansion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(partitionedExpansion)
}

func (wls *WarcraftLogs) CharacterParses(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(parses.ToDTO())
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
