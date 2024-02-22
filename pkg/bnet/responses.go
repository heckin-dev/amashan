package bnet

type CheckTokenResponse struct {
	UserName string   `json:"user_name"`
	Scope    []string `json:"scope"`
	Exp      int      `json:"exp"`
}
