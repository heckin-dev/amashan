package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var CharacterContextKey = "character"

type Character struct{}

func (r *Character) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		character, ok := vars["character"]
		if !ok {
			http.Error(w, "character not provided in route parameter", http.StatusBadRequest)
			return
		}

		character = strings.ToLower(character)

		if len(character) < 2 || len(character) > 12 {
			http.Error(w, "character name must be between 2-12 characters", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), CharacterContextKey, character)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseCharacter() *Character {
	return &Character{}
}
