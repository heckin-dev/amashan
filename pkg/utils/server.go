package utils

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// StartServerWithGracefulShutdown takes the provides mux and bind address and starts the server
// with a graceful shutdown. This shutdown will block for 30 seconds in an attempt to let other
// tasks have time to finish.
func StartServerWithGracefulShutdown(mux http.Handler, addr string, l hclog.Logger) {
	// Server configuration
	srv := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Error("listener", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	l.Info("Terminating", "signal", <-sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("shutdown", "error", err)
	}
	l.Info("graceful shutdown complete")
}
