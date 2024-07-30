package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/wl"
	"net/http"
)

type WarcraftLogs struct {
	l hclog.Logger

	client *wl.WarcraftLogsClient
}

func (w *WarcraftLogs) Route(r *mux.Router) {
	wlRouter := r.PathPrefix("/warcraftlogs").Subrouter()

	wlRouter.HandleFunc("", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, WarcraftLogs"))
	})
}

func NewWarcraftLogs(l hclog.Logger) *WarcraftLogs {
	return &WarcraftLogs{
		l:      l,
		client: wl.NewWarcraftLogsClient(l),
	}
}
