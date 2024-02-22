package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/handlers/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockBattlenet() (*BattleNet, *httptest.Server) {
	// Mock server & token exchange
	sm := mux.NewRouter()
	mock.NewOAuth2Mock().Route(sm)
	srv := httptest.NewServer(sm)

	bnet := NewBattleNet(hclog.Default())
	bnet.client.SetConfig(&oauth2.Config{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("http://%s/authorize", srv.Listener.Addr()),
			TokenURL: fmt.Sprintf("http://%s/token", srv.Listener.Addr()),
		},
		RedirectURL: "http://localhost:9090/api/auth/battlenet/callback",
		Scopes:      []string{"wow.profile", "openid"},
	})

	return bnet, srv
}

func TestBattleNet_Authorize(t *testing.T) {
	bnet, srv := mockBattlenet()
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, "/auth/battlenet", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bnet.Authorize)

	handler.ServeHTTP(rr, req)

	got := rr.Code
	want := http.StatusTemporaryRedirect

	assert.Truef(t, got == want, "got status %v, wanted %v", got, want)

	cookies := rr.Result().Cookies()

	assert.Truef(t, len(cookies) == 1, "expected 1 cookie, got none")
	assert.Truef(t, strings.EqualFold(cookies[0].Name, "oauth"), "expcted cookie named 'oauth', got %v", cookies[0].Name)
}

func TestBattleNet_Callback(t *testing.T) {
	type args struct {
		state       string
		code        string
		closeServer bool
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Should 200",
			args: args{
				state:       "123",
				code:        "abcdefg",
				closeServer: false,
			},
			want: http.StatusOK,
		},
		{
			name: "Should 400",
			args: args{
				state:       "456",
				code:        "abcdef",
				closeServer: false,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Should 500",
			args: args{
				state:       "123",
				code:        "abcdef",
				closeServer: true,
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bnet, srv := mockBattlenet()
			defer srv.Close()

			if tt.args.closeServer {
				srv.Close()
			}

			url := fmt.Sprintf("/auth/battlenet/callback?state=%s&code=%s", tt.args.state, tt.args.code)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			// Store a session for the request.
			session, _ := bnet.store.Get(req, "oauth")
			session.Values["state"] = "123"
			session.IsNew = false
			_ = session.Save(req, rr)

			handler := http.HandlerFunc(bnet.Callback)
			handler.ServeHTTP(rr, req)

			assert.Truef(t, tt.want == rr.Code, "got status %v, wanted %v", rr.Code, tt.want)
		})
	}
}
