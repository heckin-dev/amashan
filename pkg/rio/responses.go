package rio

import "time"

type CharacterProfileResponse struct {
	CharacterInfo
	MythicPlusBestRuns       []MythicPlusRun         `json:"mythic_plus_best_runs"`
	MythicPlusRecentRuns     []MythicPlusRun         `json:"mythic_plus_recent_runs"`
	MythicPlusScoresBySeason []MythicPlusSeasonScore `json:"mythic_plus_scores_by_season"`
	MythicPlusRanks          MythicPlusRanks         `json:"mythic_plus_ranks"`
}

type CharacterInfo struct {
	Name           string    `json:"name"`
	Race           string    `json:"race"`
	Class          string    `json:"class"`
	ActiveSpecName string    `json:"active_spec_name"`
	ActiveSpecRole string    `json:"active_spec_role"`
	HonorableKills int       `json:"honorable_kills"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Region         string    `json:"region"`
	Realm          string    `json:"realm"`
	LastCrawledAt  time.Time `json:"last_crawled_at"`
	ProfileURL     string    `json:"profile_url"`
}

type MythicPlusRun struct {
	Dungeon             string            `json:"dungeon"`
	ShortName           string            `json:"short_name"`
	MythicLevel         int               `json:"mythic_level"`
	CompletedAt         time.Time         `json:"completed_at"`
	ClearTimeMS         uint64            `json:"clear_time_ms"`
	ParTimeMS           uint64            `json:"par_time_ms"`
	NumKeystoneUpgrades int               `json:"num_keystone_upgrades"`
	MapChallengeModeID  int               `json:"map_challenge_mode_id"`
	ZoneID              int               `json:"zone_id"`
	Score               float64           `json:"score"`
	Affixes             []MythicPlusAffix `json:"affixes"`
	URL                 string            `json:"url"`
}

type MythicPlusAffix struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	WowheadURL  string `json:"wowhead_url"`
}

type MythicPlusSeasonScore struct {
	Season   string   `json:"season"`
	Scores   Scores   `json:"scores"`
	Segments Segments `json:"segments"`
}

type Scores struct {
	All    float64 `json:"all"`
	DPS    float64 `json:"dps"`
	Healer float64 `json:"healer"`
	Tank   float64 `json:"tank"`
	Spec0  float64 `json:"spec_0"`
	Spec1  float64 `json:"spec_1"`
	Spec2  float64 `json:"spec_2"`
	Spec3  float64 `json:"spec_3"`
}

type Segments struct {
	All    ScoreAndColor `json:"all"`
	DPS    ScoreAndColor `json:"dps"`
	Healer ScoreAndColor `json:"healer"`
	Tank   ScoreAndColor `json:"tank"`
	Spec0  ScoreAndColor `json:"spec_0"`
	Spec1  ScoreAndColor `json:"spec_1"`
	Spec2  ScoreAndColor `json:"spec_2"`
	Spec3  ScoreAndColor `json:"spec_3"`
}

type ScoreAndColor struct {
	Score float64 `json:"score"`
	Color string  `json:"color"`
}

type MythicPlusRanks struct {
	Overall WorldRegionAndRealm `json:"overall"`
	Class   WorldRegionAndRealm `json:"class"`
	// TODO: There are other Spec/DPS (/TANK/HEALER?) fields here. The Spec are spec_* which will require additional
	// 		unmarshalling, for now we'll just ignore it.
}

type WorldRegionAndRealm struct {
	World  int `json:"world"`
	Region int `json:"region"`
	Realm  int `json:"realm"`
}
