package wl

import (
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

func (e Encounter) DTO() EncounterDTO {
	return EncounterDTO{
		ID:   int(e.ID),
		Name: string(e.Name),
	}
}

type Partition struct {
	ID          graphql.Int
	Name        graphql.String
	CompactName graphql.String
	Default     graphql.Boolean
}

type CharacterData struct {
	Character *Character `graphql:"character(name: $name, serverSlug: $server_slug, serverRegion: $server_region)"`
}

type Character struct {
	Hidden       *graphql.Boolean
	ZoneRankings *ZoneRanking `graphql:"zoneRankings(zoneID: $zone_id, partition: $partition)" scalar:"true"`
}

type ZoneRanking struct {
	BestPerformanceAverage   *float64   `json:"bestPerformanceAverage"`
	MedianPerformanceAverage *float64   `json:"medianPerformanceAverage"`
	Difficulty               *int       `json:"difficulty"`
	Metric                   *string    `json:"metric"`
	Partition                *int       `json:"partition"`
	Zone                     *int       `json:"zone"`
	AllStars                 []*AllStar `json:"allStars"`
	Rankings                 []*Ranking `json:"rankings"`
}

func (r ZoneRanking) DTO() *ZoneRankingDTO {
	var allStars []*AllStarDTO
	for _, as := range r.AllStars {
		allStars = append(allStars, as.DTO())
	}

	var rankings []*RankingDTO
	for _, r := range r.Rankings {
		rankings = append(rankings, r.DTO())
	}

	return &ZoneRankingDTO{
		BestPerformanceAverage:   r.BestPerformanceAverage,
		MedianPerformanceAverage: r.MedianPerformanceAverage,
		Difficulty:               r.Difficulty,
		Metric:                   r.Metric,
		Partition:                r.Partition,
		Zone:                     r.Zone,
		AllStars:                 allStars,
		Rankings:                 rankings,
	}
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

func (s *AllStar) DTO() *AllStarDTO {
	return &AllStarDTO{
		Partition:      s.Partition,
		Spec:           s.Spec,
		Points:         s.Points,
		PossiblePoints: s.PossiblePoints,
		Rank:           s.Rank,
		RegionRank:     s.RegionRank,
		ServerRank:     s.ServerRank,
		RankPercent:    s.RankPercent,
		Total:          s.Total,
	}
}

type Ranking struct {
	Encounter     Encounter `json:"encounter"`
	RankPercent   *float64  `json:"rankPercent"`
	MedianPercent *float64  `json:"medianPercent"`
	LockedIn      bool      `json:"lockedIn"`
	TotalKills    int       `json:"totalKills"`
	FastestKill   *uint64   `json:"fastestKill"`
	AllStars      *AllStar  `json:"allStars"`
	Spec          *string   `json:"spec"`
	BestSpec      *string   `json:"bestSpec"`
	BestAmount    float64   `json:"bestAmount"`
}

func (r Ranking) DTO() *RankingDTO {
	dto := &RankingDTO{
		Encounter:     r.Encounter.DTO(),
		RankPercent:   r.RankPercent,
		MedianPercent: r.MedianPercent,
		LockedIn:      r.LockedIn,
		TotalKills:    r.TotalKills,
		FastestKill:   r.FastestKill,
		Spec:          r.Spec,
		BestSpec:      r.BestSpec,
		BestAmount:    r.BestAmount,
	}

	if r.AllStars != nil {
		dto.AllStars = r.AllStars.DTO()
	}

	return dto
}
