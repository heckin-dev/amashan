package bnet

import (
	"context"
	"github.com/heckin-dev/amashan/pkg/middleware"
	"golang.org/x/oauth2"
	"io"
)

type AccountSummaryOptions struct {
	Token  *oauth2.Token
	Region string
}

type CharacterOptions struct {
	Region    string
	Realm     string
	Character string
}

var RegionsMap = map[string]RegionOption{
	"us": RegionUS,
	"eu": RegionEU,
	"kr": RegionKR,
	"tw": RegionTW,
}

type RegionOption string

func (r RegionOption) String() string {
	return string(r)
}

const (
	RegionUS RegionOption = "us"
	RegionEU              = "eu"
	RegionKR              = "kr"
	RegionTW              = "tw"
)

type Namespace string

const (
	ProfileNamespace Namespace = "profile"
	DynamicNamespace           = "dynamic"
	StaticNamespace            = "static"
)

type RequestOptions struct {
	Region      string
	Namespace   Namespace
	Endpoint    string
	Method      string
	Body        io.Reader
	QueryParams map[string]string
}

type MythicSeasonOptions struct {
	CharacterOptions
	Season int
}

// CharacterOptionsFromContext creates a CharacterOption from a given context.
//
// It is expected that the context contains the middleware.RegionContextKey, middleware.RealmContextKey &
// middleware.CharacterContextKey
func CharacterOptionsFromContext(ctx context.Context) *CharacterOptions {
	return &CharacterOptions{
		Region:    ctx.Value(middleware.RegionContextKey).(string),
		Realm:     ctx.Value(middleware.RealmContextKey).(string),
		Character: ctx.Value(middleware.CharacterContextKey).(string),
	}
}
