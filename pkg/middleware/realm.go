package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var RealmContextKey = "realm"

type Realm struct{}

func (r *Realm) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		realm, ok := vars[RealmContextKey]
		if !ok {
			http.Error(w, "realm not provided in route parameter", http.StatusBadRequest)
			return
		}

		realm = strings.ToLower(realm)

		ctx := context.WithValue(r.Context(), RealmContextKey, realm)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseRealm() *Realm {
	return &Realm{}
}
