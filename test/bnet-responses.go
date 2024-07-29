package test

import _ "embed"

//go:embed character-summary.json
var CharacterSummary []byte

//go:embed character-status.json
var CharacterStatus []byte

//go:embed character-equipment.json
var CharacterEquipment []byte

//go:embed character-media.json
var CharacterMedia []byte

//go:embed character-statistics.json
var CharacterStatistics []byte

//go:embed character-dungeon-encounters.json
var CharacterDungeonEncounters []byte

//go:embed character-raid-encounters.json
var CharacterRaidEncounters []byte

//go:embed mythic-keystone-index.json
var MythicKeystoneIndex []byte

//go:embed mythic-keystone-season.json
var MythicKeystoneSeason []byte
