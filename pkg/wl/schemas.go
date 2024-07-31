package wl

import (
	"encoding/json"
	"github.com/hasura/go-graphql-client"
)

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

type CharacterData struct {
	Character Character `graphql:"character(name: $name, serverSlug: $server_slug, serverRegion: $server_region)"`
}

type Character struct {
	Hidden       graphql.Boolean
	ZoneRankings ZoneRankings `graphql:"zoneRankings(zoneID: $zone_id, partition: $partition)" scalar:"true"`
}

type ZoneRankings []ZoneRanking

func (z *ZoneRankings) UnmarshalJSON(bytes []byte) error {
	var zr ZoneRanking

	err := json.Unmarshal(bytes, &zr)
	if err != nil {
		return err
	}

	*z = append(*z, zr)

	return nil
}

type ZoneRanking struct {
	BestPerformanceAverage   float64   `json:"bestPerformanceAverage"`
	MedianPerformanceAverage float64   `json:"medianPerformanceAverage"`
	Difficulty               int       `json:"difficulty"`
	Metric                   string    `json:"metric"`
	Partition                int       `json:"partition"`
	Zone                     int       `json:"zone"`
	AllStars                 []AllStar `json:"all_stars"`
	Rankings                 []Ranking `json:"rankings"`
}

type AllStar struct {
	Partition      int     `json:"partition"`
	Spec           string  `json:"spec"`
	Points         float64 `json:"points"`
	PossiblePoints float64 `json:"possiblePoints"`
	Rank           int     `json:"rank"`
	RegionRank     int     `json:"regionRank"`
	ServerRank     int     `json:"serverRank"`
	RankPercent    float64 `json:"rankPercent"`
	Total          float64 `json:"total"`
}

type Ranking struct {
	Encounter     Encounter `json:"encounter"`
	RankPercent   float64   `json:"rankPercent"`
	MedianPercent float64   `json:"medianPercent"`
	LockedIn      bool      `json:"lockedIn"`
	TotalKills    int       `json:"totalKills"`
	FastestKill   uint64    `json:"fastestKill"`
	AllStars      AllStar   `json:"allStars"`
	Spec          string    `json:"spec"`
	BestSpec      string    `json:"bestSpec"`
	BestAmount    float64   `json:"bestAmount"`
}
