package middleware

import (
	"github.com/hashicorp/go-hclog"
	"net/http"
	"time"
)

// Logging is a middleware handler for request logging.
type Logging struct {
	l hclog.Logger
}

func (l *Logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.l.Info("Request", "raddr", r.RemoteAddr, "method", r.Method, "path", r.URL.Path, "took", time.Since(start))
	})
}

// UseLogging constructs a new Logging middleware handler
func UseLogging(l hclog.Logger) *Logging {
	return &Logging{
		l: l,
	}
}
