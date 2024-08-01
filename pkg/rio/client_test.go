package rio

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

// newMockedClient creates a mocked RaiderIOClient for the created httptest.Server
func newMockedClient() (*RaiderIOClient, *httptest.Server) {
	sm := mux.NewRouter()
	NewRaiderIOMock().Route(sm)
	srv := httptest.NewServer(sm)

	client := NewRaiderIOClient(hclog.Default())
	client.apiURLFn = func() string {
		return fmt.Sprintf("http://%s", srv.Listener.Addr())
	}

	return client, srv
}

func TestRaiderIOClient_CharacterProfile(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterProfileOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterProfileResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterProfileOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "skkzr",
				},
			},
			want: &CharacterProfileResponse{
				CharacterInfo: CharacterInfo{
					Name:   "Skkzr",
					Region: "us",
					Realm:  "Illidan",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, srv := newMockedClient()
			defer srv.Close()

			got, err := r.CharacterProfile(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Realm, got.Realm)
			assert.Equal(t, tt.want.Region, got.Region)
			assert.NotEmpty(t, got.MythicPlusBestRuns)
			assert.NotEmpty(t, got.MythicPlusRecentRuns)
			assert.NotEmpty(t, got.MythicPlusBestRuns)
			assert.NotEmpty(t, got.MythicPlusRanks)
		})
	}
}
