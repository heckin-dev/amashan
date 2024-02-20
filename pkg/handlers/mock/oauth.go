package mock

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type OAuth2Mock struct{}

func (o *OAuth2Mock) Token(w http.ResponseWriter, r *http.Request) {
	token := oauth2.Token{
		AccessToken: "abcdef-ghijkl",
		TokenType:   "bearer",
		Expiry:      time.Now().Add(time.Hour * 1),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

func (o *OAuth2Mock) Route(r *mux.Router) {
	r.HandleFunc("/token", o.Token).Methods(http.MethodPost)
}

func NewOAuth2Mock() *OAuth2Mock {
	return &OAuth2Mock{}
}
