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
	Character          Link                `json:"character"`
	ProtectedCharacter Link                `json:"protected_character"`
	Name               string              `json:"name"`
	ID                 int                 `json:"id"`
	Realm              Realm               `json:"realm"`
	PlayableClass      PlayableRaceOrClass `json:"playable_class"`
	PlayableRace       PlayableRaceOrClass `json:"playable_race"`
	Gender             GenderOrFaction     `json:"gender"`
	Faction            GenderOrFaction     `json:"faction"`
	Level              int                 `json:"level"`
}

type Realm struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type PlayableRaceOrClass struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GenderOrFaction struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Link struct {
	Href string `json:"href"`
}
