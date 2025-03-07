package bnet

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/heckin-dev/amashan/pkg/handlers/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
	"net/http/httptest"
	"os"
	"testing"
)

func newMockedClient() (*BattlenetClient, *httptest.Server) {
	// Mock server & token exchange
	sm := mux.NewRouter()
	mock.NewOAuth2Mock().Route(sm)
	NewBattleNetMock().Route(sm)
	srv := httptest.NewServer(sm)

	os.Setenv("SESSION_KEY", "catswithhats")

	bnet := NewBattlnetClient(hclog.Default())

	bnet.clientConfig = &cc.Config{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		TokenURL:     fmt.Sprintf("http://%s/token", srv.Listener.Addr()),
	}

	bnet.oauthConfig = &oauth2.Config{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("http://%s/authorize", srv.Listener.Addr()),
			TokenURL: fmt.Sprintf("http://%s/token", srv.Listener.Addr()),
		},
		RedirectURL: "http://localhost:9090/api/auth/battlenet/callback",
		Scopes:      []string{"wow.profile", "openid"},
	}

	bnet.apiURLFn = func(region string) string {
		return fmt.Sprintf("http://%s", srv.Listener.Addr())
	}

	return bnet, srv
}

func TestBattlenetClient_CharacterEquipmentSummary(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterEquipmentResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "tichondrius",
					Character: "skxyz",
				},
			},
			want: &CharacterEquipmentResponse{
				Character: Character{
					Name: "Skxyz",
					ID:   242328372,
					Realm: Realm{
						Slug: "tichondrius",
					},
				},
			},
			wantErr: false,
		},
		//{
		//	name: "Should 400 - Missing Realm",
		//	args: args{
		//		ctx: nil,
		//		options: &CharacterOptions{
		//			Region:    "us",
		//			Realm:     "",
		//			Character: "aulene",
		//		},
		//	},
		//	want:    nil,
		//	wantErr: true,
		//},
		//{
		//	name: "Should 400 - Missing Character",
		//	args: args{
		//		ctx: nil,
		//		options: &CharacterOptions{
		//			Region:    "us",
		//			Realm:     "illidan",
		//			Character: "",
		//		},
		//	},
		//	want:    nil,
		//	wantErr: true,
		//},
		{
			name: "Should 400 - Character too short",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "a",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Character too long",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "areallylongusernamethatisntreal",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterEquipmentSummary(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterEquipmentSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
			assert.NotEmpty(t, got.EquippedItems)
			assert.NotEmpty(t, got.EquippedItemSets)
		})
	}
}

func TestBattlenetClient_CharacterMedia(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterMediaResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "aulene",
				},
			},
			want: &CharacterMediaResponse{
				Character: Character{
					Name: "Aulene",
					ID:   229483897,
					Realm: Realm{
						Slug: "illidan",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should 400 - Missing Realm",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "",
					Character: "aulene",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Missing Character",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Character too short",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "a",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Character too long",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "areallylongusernamethatisntreal",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterMedia(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
			assert.NotEmpty(t, got.Assets)
		})
	}
}

func TestBattlenetClient_MythicKeystoneIndex(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *MythicKeystoneIndexResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Skkzr",
				},
			},
			want: &MythicKeystoneIndexResponse{
				Character: Character{
					Name: "Skkzr",
					ID:   225511351,
					Realm: Realm{
						Slug: "illidan",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should 400 - Missing Realm",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "",
					Character: "aulene",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Missing Character",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Character too short",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "a",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should 400 - Character too long",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "areallylongusernamethatisntreal",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.MythicKeystoneIndex(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("MythicKeystoneIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
			assert.NotEmpty(t, got.CurrentMythicRating)
			assert.NotEmpty(t, got.Seasons)
			assert.NotEmpty(t, got.CurrentPeriod)
		})
	}
}

func TestBattlenetClient_MythicKeystoneSeason(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *MythicSeasonOptions
	}
	tests := []struct {
		name      string
		args      args
		want      *MythicKeystoneSeasonResponse
		wantErr   bool
		should404 bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &MythicSeasonOptions{
					CharacterOptions: CharacterOptions{
						Region:    "us",
						Realm:     "illidan",
						Character: "Skkzr",
					},
					Season: 12,
				},
			},
			want: &MythicKeystoneSeasonResponse{
				Character: &Character{
					Name: "Skkzr",
					ID:   225511351,
					Realm: Realm{
						Slug: "illidan",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should 404",
			args: args{
				ctx: nil,
				options: &MythicSeasonOptions{
					CharacterOptions: CharacterOptions{
						Region:    "us",
						Realm:     "",
						Character: "aulene",
					},
					Season: 12,
				},
			},
			want:      nil,
			wantErr:   false,
			should404: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.MythicKeystoneSeason(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("MythicKeystoneSeason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if tt.should404 {
				assert.False(t, got.CharacterPlayedSeason)
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
			assert.NotEmpty(t, got.MythicRating)
			assert.NotEmpty(t, got.BestRuns)
			assert.NotEmpty(t, got.Season)
			assert.True(t, got.CharacterPlayedSeason)
		})
	}
}

func TestBattlenetClient_CharacterStatistics(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterStatisticsResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Aulene",
				},
			},
			want: &CharacterStatisticsResponse{
				Character: Character{
					Name: "Aulene",
					ID:   229483897,
					Realm: Realm{
						Slug: "illidan",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterStatistics(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
		})
	}
}

func TestBattlenetClient_CharacterDungeonEncounters(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Aulene",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterDungeonEncounters(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterDungeonEncounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.NotEmpty(t, got.Expansions)
		})
	}
}

func TestBattlenetClient_CharacterRaidEncounters(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterRaidEncountersResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Skkzr",
				},
			},
			want: &CharacterRaidEncountersResponse{
				Character: Character{
					Name: "Skkzr",
					ID:   225511351,
					Realm: Realm{
						Slug: "illidan",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterRaidEncounters(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterRaidEncounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Character.Name, got.Character.Name)
			assert.Equal(t, tt.want.Character.ID, got.Character.ID)
			assert.Equal(t, tt.want.Character.Realm.Slug, got.Character.Realm.Slug)
			assert.NotEmpty(t, got.Expansions)
		})
	}
}

func TestBattlenetClient_CharacterSummary(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterSummaryResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Aulene",
				},
			},
			want: &CharacterSummaryResponse{
				Name: "Aulene",
				ID:   229483897,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterSummary(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.ID, got.ID)
		})
	}
}

func TestBattlenetClient_CharacterStatus(t *testing.T) {
	type args struct {
		ctx     context.Context
		options *CharacterOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *CharacterStatusResponse
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx: nil,
				options: &CharacterOptions{
					Region:    "us",
					Realm:     "illidan",
					Character: "Aulene",
				},
			},
			want: &CharacterStatusResponse{
				ID:      229483897,
				IsValid: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.CharacterStatus(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CharacterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.IsValid, got.IsValid)
			assert.Equal(t, tt.want.ID, got.ID)
		})
	}
}

func TestBattlenetClient_RealmsByRegion(t *testing.T) {
	type args struct {
		ctx    context.Context
		option RegionOption
	}
	tests := []struct {
		name    string
		args    args
		want    Realm
		wantErr bool
	}{
		{
			name: "Should 200",
			args: args{
				ctx:    nil,
				option: RegionsMap["kr"],
			},
			want: Realm{
				NamedTypeAndID: NamedTypeAndID{
					KeyedID: KeyedID{
						Key: Link{
							Href: "https://kr.api.blizzard.com/data/wow/realm/205?namespace=dynamic-kr",
						},
						ID: 205,
					},
					Name: toString("Azshara"),
				},
				Slug: "azshara",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, srv := newMockedClient()
			defer srv.Close()

			got, err := b.RealmsByRegion(tt.args.ctx, tt.args.option)
			if (err != nil) != tt.wantErr {
				t.Errorf("RealmsByRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			assert.Contains(t, got.Realms, tt.want)
		})
	}
}

func toString(v string) *string {
	return &v
}
