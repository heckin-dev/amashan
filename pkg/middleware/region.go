package middleware

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"slices"
	"strings"
)

var regions = []string{"us", "eu", "kr", "tw"}

var RegionContextKey = "region"

type Region struct{}

func (r *Region) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		region, ok := vars["region"]
		if !ok {
			http.Error(w, "region not provided in route parameter", http.StatusBadRequest)
			return
		}

		region = strings.ToLower(region)

		if !slices.Contains(regions, region) {
			http.Error(w, fmt.Sprintf("region '%s' is not a supported region", region), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), RegionContextKey, region)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseRegion() *Region {
	return &Region{}
}
