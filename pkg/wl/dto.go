package wl

type CharacterParseDTO struct {
	Hidden       *bool           `json:"hidden,omitempty"`
	ZoneRankings *ZoneRankingDTO `json:"zone_rankings,omitempty"`
}

type ZoneRankingDTO struct {
	BestPerformanceAverage   *float64      `json:"best_performance_average,omitempty"`
	MedianPerformanceAverage *float64      `json:"median_performance_average,omitempty"`
	Difficulty               *int          `json:"difficulty,omitempty"`
	Metric                   *string       `json:"metric,omitempty"`
	Partition                *int          `json:"partition,omitempty"`
	Zone                     *int          `json:"zone,omitempty"`
	AllStars                 []*AllStarDTO `json:"all_stars,omitempty"`
	Rankings                 []*RankingDTO `json:"rankings,omitempty"`
}

type AllStarDTO struct {
	Partition      int     `json:"partition"`
	Spec           string  `json:"spec"`
	Points         float64 `json:"points"`
	PossiblePoints float64 `json:"possible_points"`
	Rank           int     `json:"rank"`
	RegionRank     int     `json:"region_rank"`
	ServerRank     int     `json:"server_rank"`
	RankPercent    float64 `json:"rank_percent"`
	Total          float64 `json:"total"`
}

type RankingDTO struct {
	Encounter     EncounterDTO `json:"encounter"`
	RankPercent   *float64     `json:"rank_percent,omitempty"`
	MedianPercent *float64     `json:"median_percent,omitempty"`
	LockedIn      bool         `json:"locked_in"`
	TotalKills    int          `json:"total_kills"`
	FastestKill   *uint64      `json:"fastest_kill,omitempty"`
	AllStars      *AllStarDTO  `json:"all_stars,omitempty"`
	Spec          *string      `json:"spec,omitempty"`
	BestSpec      *string      `json:"best_spec,omitempty"`
	BestAmount    float64      `json:"best_amount"`
}

type EncounterDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PartitionedExpansion struct {
	ID    int             `json:"id"`
	Name  string          `json:"name"`
	Zones []ExpansionZone `json:"zones"`
}

func PartitionedExpansionFromExpansionEncounterQuery(q *ExpansionEncountersQuery) []*PartitionedExpansion {
	var result []*PartitionedExpansion

	for _, exp := range q.WorldData.Expansions {
		pe := &PartitionedExpansion{
			ID:    int(exp.ID),
			Name:  string(exp.Name),
			Zones: nil,
		}

		// Collects the zones
		var zones []ExpansionZone
		for _, zone := range exp.Zones {
			peZone := ExpansionZone{
				ID:         int(zone.ID),
				Name:       string(zone.Name),
				Encounters: zone.GetZoneEncounters(),
				Partitions: zone.GetZonePartitions(),
			}
			zones = append(zones, peZone)
		}
		pe.Zones = zones
		result = append(result, pe)
	}

	return result
}

type ExpansionZone struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Encounters []ZoneEncounter `json:"encounters"`
	Partitions []ZonePartition `json:"partitions"`
}

type ZoneEncounter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ZonePartition struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CompactName string `json:"compact_name"`
	Default     bool   `json:"default"`
}
