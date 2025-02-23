package config

import "time"

// Config is a struct holds the client configuration options
type Config struct {
	RelyingPartyName string
	RelyingPartyUUID string
	HashType         string
	Text             string
	TextFormat       string
	Language         string
	URL              string
	Timeout          time.Duration
}
