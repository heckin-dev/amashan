package bnet

// CheckTokenResponse /oauth/check_token
type CheckTokenResponse struct {
	UserName string   `json:"user_name"`
	Scope    []string `json:"scope"`
}

// UserInfoResponse /oauth/userinfo
type UserInfoResponse struct {
	ID        int    `json:"id"`
	BattleTag string `json:"battletag"`
}

// AccountSummaryResponse /profile/user/wow
type AccountSummaryResponse struct {
	WowAccounts []AccountSummaryAccount `json:"wow_accounts"`
}

type AccountSummaryAccount struct {
	ID         int                       `json:"id"`
	Characters []AccountSummaryCharacter `json:"characters"`
}

type AccountSummaryCharacter struct {
	Character          Link           `json:"character"`
	ProtectedCharacter Link           `json:"protected_character"`
	Name               string         `json:"name"`
	ID                 int            `json:"id"`
	Realm              Realm          `json:"realm"`
	PlayableClass      NamedTypeAndID `json:"playable_class"`
	PlayableRace       NamedTypeAndID `json:"playable_race"`
	Gender             TypeAndName    `json:"gender"`
	Faction            TypeAndName    `json:"faction"`
	Level              int            `json:"level"`
}

// CharacterSummaryResponse /profile/wow/character/{realmSlug}/{characterName}
type CharacterSummaryResponse struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name"`
	Gender             TypeAndName    `json:"gender"`
	Faction            TypeAndName    `json:"faction"`
	Race               NamedTypeAndID `json:"race"`
	CharacterClass     NamedTypeAndID `json:"character_class"`
	ActiveSpec         NamedTypeAndID `json:"active_spec"`
	Realm              Realm          `json:"realm"`
	Guild              Guild          `json:"guild"`
	Level              int            `json:"level"`
	Experience         int            `json:"experience"`
	AchievementPoints  int            `json:"achievement_points"`
	LastLoginTimestamp uint64         `json:"last_login_timestamp"`
	AverageItemLevel   int            `json:"average_item_level"`
	EquippedItemLevel  int            `json:"equipped_item_level"`
}

// CharacterStatusResponse /profile/wow/character/{realmSlug}/{characterName}/status
type CharacterStatusResponse struct {
	ID      int  `json:"id"`
	IsValid bool `json:"is_valid"`
}

// CharacterEquipmentResponse /profile/wow/character/{realmSlug}/{characterName}/equipment
type CharacterEquipmentResponse struct {
	Character        Character                           `json:"character"`
	EquippedItems    []CharacterEquipmentEquippedItem    `json:"equipped_items"`
	EquippedItemSets []CharacterEquipmentEquippedItemSet `json:"equipped_item_sets"`
}

type CharacterEquipmentEquippedItem struct {
	Item                 KeyedID            `json:"item"`
	Enchantments         []*ItemEnchantment `json:"enchantments,omitempty"`
	Sockets              []*ItemSocket      `json:"sockets"`
	Slot                 TypeAndName        `json:"slot"`
	Quantity             int                `json:"quantity"`
	Context              int                `json:"context"`
	BonusList            []int              `json:"bonus_list"`
	Quality              TypeAndName        `json:"quality"`
	Name                 string             `json:"name"`
	ModifiedAppearanceID *int               `json:"modified_appearance_id,omitempty"`
	Media                KeyedID            `json:"media"`
	ItemClass            NamedTypeAndID     `json:"item_class"`
	ItemSubclass         NamedTypeAndID     `json:"item_subclass"`
	InventoryType        TypeAndName        `json:"inventory_type"`
	Binding              TypeAndName        `json:"binding"`
	Armor                *Armor             `json:"armor,omitempty"`
	Weapon               *Armor             `json:"weapon,omitempty"`
	Stats                []ItemStat         `json:"stats"`
	SellPrice            SellPrice          `json:"sell_price"`
	Requirements         ItemRequirements   `json:"requirements"`
	Set                  *Set               `json:"set,omitempty"`
	Level                IntDisplayValue    `json:"level"`
	Transmog             *Transmog          `json:"transmog,omitempty"`
	Durability           *IntDisplayValue   `json:"durability,omitempty"`
	NameDescription      *Display           `json:"name_description,omitempty"`
	IsSubclassHidden     *bool              `json:"is_subclass_hidden,omitempty"`
}

type ItemEnchantment struct {
	DisplayString   string         `json:"display_string"`
	SourceItem      NamedTypeAndID `json:"source_item"`
	EnchantmentID   int            `json:"enchantment_id"`
	EnchantmentSlot TypeAndID      `json:"enchantment_slot"`
}

type ItemSocket struct {
	SocketType    TypeAndName    `json:"socket_type"`
	Item          NamedTypeAndID `json:"item"`
	DisplayString string         `json:"display_string"`
	Media         KeyedID        `json:"media"`
}

type CharacterEquipmentEquippedItemSet struct {
	ItemSet       NamedTypeAndID `json:"item_set"`
	Items         []SetItem      `json:"items"`
	Effects       []SetEffect    `json:"effects"`
	DisplayString string         `json:"display_string"`
}

// CharacterMediaResponse /profile/wow/character/{realmSlug}/{characterName}/character-media
type CharacterMediaResponse struct {
	Character Character     `json:"character"`
	Assets    []KeyAndValue `json:"assets"`
}

// MythicKeystoneIndexResponse /profile/wow/character/{realmSlug}/{characterName}/mythic-keystone-profile
type MythicKeystoneIndexResponse struct {
	Character           Character           `json:"character"`
	Seasons             []KeyedID           `json:"seasons"`
	CurrentPeriod       CurrentMythicPeriod `json:"current_period"`
	CurrentMythicRating MythicRating        `json:"current_mythic_rating"`
}

// MythicKeystoneSeasonResponse /profile/wow/character/{realmSlug}/{characterName}/mythic-keystone-profile/season/{seasonId}
type MythicKeystoneSeasonResponse struct {
	CharacterPlayedSeason bool          `json:"character_played_season"`
	Character             *Character    `json:"character,omitempty"`
	Season                *KeyedID      `json:"season,omitempty"`
	BestRuns              []*MythicRun  `json:"best_runs,omitempty"`
	MythicRating          *MythicRating `json:"mythic_rating,omitempty"`
}

type CurrentMythicPeriod struct {
	Period   KeyedID     `json:"period"`
	BestRuns []MythicRun `json:"best_runs"`
}

type MythicRun struct {
	CompletedTimestamp    uint64                 `json:"completed_timestamp"`
	Duration              uint64                 `json:"duration"`
	KeystoneLevel         int                    `json:"keystone_level"`
	KeystoneAffixes       []NamedTypeAndID       `json:"keystone_affixes"`
	Members               []MythicKeystoneMember `json:"members"`
	Dungeon               NamedTypeAndID         `json:"dungeon"`
	IsCompletedWithinTime bool                   `json:"is_completed_within_time"`
	MythicRating          MythicRating           `json:"mythic_rating"`
	MapRating             MythicRating           `json:"map_rating"`
}

type MythicKeystoneMember struct {
	Character         Character      `json:"character"`
	Specialization    NamedTypeAndID `json:"specialization"`
	Race              NamedTypeAndID `json:"race"`
	EquippedItemLevel int            `json:"equipped_item_level"`
}

type MythicRating struct {
	Rating float64   `json:"rating"`
	Color  ColorRGBA `json:"color"`
}

// CharacterStatisticsResponse /profile/wow/character/{realmSlug}/{characterName}/statistics
type CharacterStatisticsResponse struct {
	Character                   Character              `json:"character"`
	PowerType                   NamedTypeAndID         `json:"power_type"`
	Health                      float64                `json:"health"`
	Power                       float64                `json:"power"`
	Avoidance                   StatisticRating        `json:"avoidance"`
	Block                       StatisticRating        `json:"block"`
	Dodge                       StatisticRating        `json:"dodge"`
	Lifesteal                   StatisticRating        `json:"lifesteal"`
	Mastery                     StatisticRating        `json:"mastery"`
	MeleeCrit                   StatisticRating        `json:"melee_crit"`
	MeleeHaste                  StatisticRating        `json:"melee_haste"`
	Parry                       StatisticRating        `json:"parry"`
	RangedCrit                  StatisticRating        `json:"ranged_crit"`
	RangedHaste                 StatisticRating        `json:"ranged_haste"`
	Speed                       StatisticRating        `json:"speed"`
	SpellCrit                   StatisticRating        `json:"spell_crit"`
	SpellHaste                  StatisticRating        `json:"spell_haste"`
	Agility                     StatisticEffectiveness `json:"agility"`
	Armor                       StatisticEffectiveness `json:"armor"`
	Intellect                   StatisticEffectiveness `json:"intellect"`
	Stamina                     StatisticEffectiveness `json:"stamina"`
	Strength                    StatisticEffectiveness `json:"strength"`
	AttackPower                 float64                `json:"attack_power"`
	BonusArmor                  float64                `json:"bonus_armor"`
	MainHandDPS                 float64                `json:"main_hand_dps"`
	MainHandDamageMax           float64                `json:"main_hand_damage_max"`
	MainHandDamageMin           float64                `json:"main_hand_damage_min"`
	MainHandSpeed               float64                `json:"main_hand_speed"`
	ManaRegen                   float64                `json:"mana_regen"`
	ManaRegenCombat             float64                `json:"mana_regen_combat"`
	OffHandDPS                  float64                `json:"off_hand_dps"`
	OffHandDamageMax            float64                `json:"off_hand_damage_max"`
	OffHandDamageMin            float64                `json:"off_hand_damage_min"`
	OffHandSpeed                float64                `json:"off_hand_speed"`
	SpellPenetration            float64                `json:"spell_penetration"`
	SpellPower                  float64                `json:"spell_power"`
	Versatility                 float64                `json:"versatility"`
	VersatilityDamageDone       float64                `json:"versatility_damage_done"`
	VersatilityDamageTakenBonus float64                `json:"versatility_damage_taken_bonus"`
	VersatilityHealingDone      float64                `json:"versatility_healing_done"`
}

type StatisticRating struct {
	Rating      float64  `json:"rating"`
	RatingBonus float64  `json:"rating_bonus"`
	Value       *float64 `json:"value,omitempty"`
}

type StatisticEffectiveness struct {
	Base      float64 `json:"base"`
	Effective float64 `json:"effective"`
}

// CharacterDungeonEncountersResponse /profile/wow/character/{realmSlug}/{characterName}/encounters/dungeons
type CharacterDungeonEncountersResponse struct {
	Expansions []EncounterExpansion `json:"expansions"`
}

// CharacterRaidEncountersResponse /profile/wow/character/{realmSlug}/{characterName}/encounters/raids
type CharacterRaidEncountersResponse struct {
	Character  Character            `json:"character"`
	Expansions []EncounterExpansion `json:"expansions"`
}

type EncounterExpansion struct {
	Expansion NamedTypeAndID      `json:"expansion"`
	Instances []EncounterInstance `json:"instances"`
}

type EncounterInstance struct {
	Instance NamedTypeAndID  `json:"instance"`
	Modes    []EncounterMode `json:"modes"`
}

type EncounterMode struct {
	Difficulty TypeAndName `json:"difficulty"`
	Status     TypeAndName `json:"status"`
}

type EncounterProgress struct {
	CompletedCount int         `json:"completed_count"`
	TotalCount     int         `json:"total_count"`
	Encounters     []Encounter `json:"encounters"`
}

type Encounter struct {
	Encounter         NamedTypeAndID `json:"encounter"`
	CompletedCount    int            `json:"completed_count"`
	LastKillTimestamp uint64         `json:"last_kill_timestamp"`
}

// RealmIndexResponse /data/wow/realm/index
type RealmIndexResponse struct {
	// Region is not  field is provided by Blizzard. We add this field for response clarity.
	Region string  `json:"region"`
	Realms []Realm `json:"realms"`
}

type KeyAndValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyedID struct {
	Key Link `json:"key"`
	ID  int  `json:"id"`
}

type TypeAndName struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type NamedTypeAndID struct {
	KeyedID
	Name *string `json:"name,omitempty"`
}

type Link struct {
	Href string `json:"href"`
}

type Character struct {
	Key   *Link  `json:"key,omitempty"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Realm Realm  `json:"realm"`
}

type Realm struct {
	NamedTypeAndID
	Slug string `json:"slug"`
}

type Armor struct {
	Value   int     `json:"value"`
	Display Display `json:"display"`
}

type Weapon struct {
	Damage      WeaponDamage      `json:"damage"`
	AttackSpeed FloatDisplayValue `json:"attack_speed"`
	DPS         FloatDisplayValue `json:"dps"`
}

type WeaponDamage struct {
	MinValue      int         `json:"min_value"`
	MaxValue      int         `json:"max_value"`
	DisplayString string      `json:"display_string"`
	DamageClass   TypeAndName `json:"damage_class"`
}

type Display struct {
	DisplayString string    `json:"display_string"`
	Color         ColorRGBA `json:"color"`
}

type ColorRGBA struct {
	R uint8   `json:"r"`
	G uint8   `json:"g"`
	B uint8   `json:"b"`
	A float64 `json:"a"`
}

type ItemStat struct {
	Type    TypeAndName `json:"type"`
	Value   int         `json:"value"`
	Display Display     `json:"display"`
}

type SellPrice struct {
	Value          int            `json:"value"`
	DisplayStrings DisplayStrings `json:"display_strings"`
}

type DisplayStrings struct {
	Header string `json:"header"`
	Gold   string `json:"gold"`
	Silver string `json:"silver"`
	Copper string `json:"copper"`
}

type ItemRequirements struct {
	Level           IntDisplayValue `json:"level"`
	PlayableClasses PlayableClasses `json:"playable_classes"`
}

type IntDisplayValue struct {
	Value         int    `json:"value"`
	DisplayString string `json:"display_string"`
}

type FloatDisplayValue struct {
	Value         float64 `json:"value"`
	DisplayString string  `json:"display_string"`
}

type PlayableClasses struct {
	Links         []NamedTypeAndID `json:"links"`
	DisplayString string           `json:"display_string"`
}

type Set struct {
	ItemSet       NamedTypeAndID `json:"item_set"`
	Items         []SetItem      `json:"items"`
	Effects       []SetEffect    `json:"effects"`
	DisplayString string         `json:"display_string"`
}

type SetItem struct {
	Item       NamedTypeAndID `json:"item"`
	IsEquipped bool           `json:"is_equipped"`
}

type SetEffect struct {
	DisplayString string `json:"display_string"`
	RequiredCount int    `json:"required_count"`
	IsActive      bool   `json:"is_active"`
}

type Transmog struct {
	Item                     NamedTypeAndID `json:"item"`
	DisplayString            string         `json:"display_string"`
	ItemModifiedAppearanceID int            `json:"item_modified_appearance_id"`
}

type Guild struct {
	NamedTypeAndID
	Realm   Realm       `json:"realm"`
	Faction TypeAndName `json:"faction"`
}

type TypeAndID struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}
