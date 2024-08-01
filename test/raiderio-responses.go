package test

import _ "embed"

//go:embed raiderio_mplus-current-season.json
var MPlusCurrentSeason []byte

//go:embed raiderio_mplus-current-ranks.json
var MPlusRanks []byte

//go:embed raiderio_mplus-recent-runs.json
var MPlusRecentRuns []byte

//go:embed raiderio_mplus-best-runs.json
var MPlusBestRuns []byte

//go:embed raiderio_all-in-one.json
var MPlusAllInOne []byte
