package models

type Client struct {
	ID           int
	ClientID     string
	HashedSecret string
	Scopes       []string
	IsActive     bool
}
