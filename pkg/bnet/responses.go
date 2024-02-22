package bnet

type CheckTokenResponse struct {
	UserName string   `json:"user_name"`
	Scope    []string `json:"scope"`
}

type UserInfoResponse struct {
	ID        int    `json:"id"`
	BattleTag string `json:"battletag"`
}
