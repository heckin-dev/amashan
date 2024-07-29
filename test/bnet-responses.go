package test

import _ "embed"

//go:embed character-equipment.json
var CharacterEquipment []byte

//go:embed character-media.json
var CharacterMedia []byte

//go:embed mythic-keystone-index.json
var MythicKeystoneIndex []byte

//go:embed mythic-keystone-season.json
var MythicKeystoneSeason []byte
