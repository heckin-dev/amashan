package wl

import "github.com/shurcooL/graphql"

type RateLimitData struct {
	LimitPerHour        graphql.Int
	PointsSpentThisHour graphql.Float
	PointsResetIn       graphql.Int
}

type WorldData struct {
	Expansions []Expansion
}

type Expansion struct {
	ID    graphql.Int
	Name  graphql.String
	Zones []Zone
}

type Zone struct {
	ID         graphql.Int
	Encounters []Encounter
}

type Encounter struct {
	ID   graphql.Int
	Name graphql.String
}
