package bnet

type RequestType string

const (
	ClientRequest RequestType = "client"
	OAuthRequest              = "oauth"
)
