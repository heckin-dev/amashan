package wl

import (
	"context"
	"github.com/heckin-dev/amashan/pkg/middleware"
)

type CharacterParsesQueryOptions struct {
	Name         string
	ServerSlug   string
	ServerRegion string
	ZoneID       int
	Partition    *int
}

// CharacterParsesQueryOptionsFromContext creates a CharacterParsesQueryOptions from a given context.
//
// It is expected that the context contains the middleware.RegionContextKey, middleware.RealmContextKey &
// middleware.CharacterContextKey
func CharacterParsesQueryOptionsFromContext(ctx context.Context) *CharacterParsesQueryOptions {
	return &CharacterParsesQueryOptions{
		ServerRegion: ctx.Value(middleware.RegionContextKey).(string),
		ServerSlug:   ctx.Value(middleware.RealmContextKey).(string),
		Name:         ctx.Value(middleware.CharacterContextKey).(string),
	}
}
