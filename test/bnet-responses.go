package test

import _ "embed"

//go:embed character-equipment.json
var CharacterEquipment []byte

//go:embed character-media.json
var CharacterMedia []byte

//go:embed character-statistics.json
var CharacterStatistics []byte

var CharacterEncounters []byte

var CharacterDungeonEncounters []byte

var CharacterRaidEncounters []byte

//go:embed mythic-keystone-index.json
var MythicKeystoneIndex []byte

//go:embed mythic-keystone-season.json
var MythicKeystoneSeason []byte
