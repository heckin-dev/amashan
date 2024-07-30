package wl

import "github.com/shurcooL/graphql"

type RateLimitData struct {
	LimitPerHour        graphql.Int
	PointsSpentThisHour graphql.Float
	PointsResetIn       graphql.Int
}
