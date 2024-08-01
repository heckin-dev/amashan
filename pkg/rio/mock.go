package rio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/heckin-dev/amashan/test"
	"net/http"
)

type RaiderIOMock struct{}

func (i *RaiderIOMock) CharacterProfile(w http.ResponseWriter, r *http.Request) {
	expectedParams := []string{"region", "realm", "name", "fields"}

	q := r.URL.Query()
	for _, expParam := range expectedParams {
		if !q.Has(expParam) {
			http.Error(w, fmt.Errorf("expected param: '%s' was missing", expParam).Error(), http.StatusBadRequest)
			return
		}
	}

	res := &CharacterProfileResponse{}
	err := json.NewDecoder(bytes.NewReader(test.MPlusAllInOne)).Decode(res)
	if err != nil {
		http.Error(w, "failed to decode test.MPlusAllInOne", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (i *RaiderIOMock) Route(r *mux.Router) {
	r.HandleFunc("/characters/profile", i.CharacterProfile)
}

func NewRaiderIOMock() *RaiderIOMock {
	return &RaiderIOMock{}
}
