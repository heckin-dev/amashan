package rio

import (
	"context"
	"github.com/heckin-dev/amashan/pkg/middleware"
)

type CharacterProfileOptions struct {
	Region    string
	Realm     string
	Character string
}

// CharacterProfileOptionsFromContext creates a CharacterProfileOptions from a given context.
//
// It is expected that the context contains the middleware.RegionContextKey, middleware.RealmContextKey &
// middleware.CharacterContextKey
func CharacterProfileOptionsFromContext(ctx context.Context) *CharacterProfileOptions {
	return &CharacterProfileOptions{
		Region:    ctx.Value(middleware.RegionContextKey).(string),
		Realm:     ctx.Value(middleware.RealmContextKey).(string),
		Character: ctx.Value(middleware.CharacterContextKey).(string),
	}
}
