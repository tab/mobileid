package mobileid

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/tab/mobileid/internal/config"
	"github.com/tab/mobileid/internal/errors"
	"github.com/tab/mobileid/internal/requests"
	"github.com/tab/mobileid/internal/utils"
)

const (
	Text       = "Enter PIN1"
	TextFormat = "GSM-7"
	Language   = "ENG"
	Timeout    = requests.Timeout
	URL        = "https://tsp.demo.sk.ee/mid-api"
)

type Client interface {
	CreateSession(ctx context.Context, phoneNumber, nationalIdentityNumber string) (*Session, error)
	FetchSession(ctx context.Context, sessionId string) (*Person, error)

	WithRelyingPartyName(name string) Client
	WithRelyingPartyUUID(id string) Client
	WithHashType(hashType string) Client
	WithText(text string) Client
	WithTextFormat(format string) Client
	WithLanguage(language string) Client
	WithURL(url string) Client
	WithTimeout(timeout time.Duration) Client
	WithTLSConfig(tlsConfig *tls.Config) Client

	Validate() error
}

type client struct {
	config *config.Config
}

func NewClient() Client {
	cfg := &config.Config{
		HashType:   utils.HashTypeSHA512,
		Text:       Text,
		TextFormat: TextFormat,
		Language:   Language,
		URL:        URL,
		Timeout:    Timeout,
	}

	return &client{
		config: cfg,
	}
}

func (c *client) WithRelyingPartyName(name string) Client {
	c.config.RelyingPartyName = name
	return c
}

func (c *client) WithRelyingPartyUUID(id string) Client {
	c.config.RelyingPartyUUID = id
	return c
}

func (c *client) WithHashType(hashType string) Client {
	c.config.HashType = hashType
	return c
}

func (c *client) WithText(text string) Client {
	c.config.Text = text
	return c
}

func (c *client) WithTextFormat(format string) Client {
	c.config.TextFormat = format
	return c
}

func (c *client) WithLanguage(language string) Client {
	c.config.Language = language
	return c
}

func (c *client) WithURL(url string) Client {
	c.config.URL = url
	return c
}

func (c *client) WithTimeout(timeout time.Duration) Client {
	c.config.Timeout = timeout
	return c
}

func (c *client) WithTLSConfig(tlsConfig *tls.Config) Client {
	c.config.TLSConfig = tlsConfig
	return c
}

func (c *client) Validate() error {
	if c.config.RelyingPartyName == "" {
		return errors.ErrMissingRelyingPartyName
	}

	if c.config.RelyingPartyUUID == "" {
		return errors.ErrMissingRelyingPartyUUID
	}

	return nil
}
