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

type WorldDataByExpansionID struct {
	Zones []PartitionedZone `graphql:"zones(expansion_id: $expansion_id)"`
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

type PartitionedZone struct {
	ID         graphql.Int
	Name       graphql.String
	Partitions []Partition
}

type Encounter struct {
	ID   graphql.Int
	Name graphql.String
}

type Partition struct {
	ID          graphql.Int
	Name        graphql.String
	CompactName graphql.String
	Default     graphql.Boolean
}
