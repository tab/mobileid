package mobileid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tab/mobileid/internal/config"
	"github.com/tab/mobileid/internal/errors"
)

func Test_NewClient(t *testing.T) {
	type result struct {
		config *config.Config
	}

	tests := []struct {
		name     string
		before   func(c Client)
		expected result
	}{
		{
			name: "Success",
			before: func(c Client) {
				c.
					WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
					WithHashType("SHA512").
					WithText("Enter PIN1").
					WithTextFormat("GSM-7").
					WithLanguage("ENG").
					WithURL("https://tsp.demo.sk.ee/mid-api").
					WithTimeout(60 * time.Second)
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
				},
			},
		},
		{
			name: "Default values",
			before: func(c Client) {
				c.
					WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000")
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
				},
			},
		},
		{
			name: "Error: Missing relying party name",
			before: func(c Client) {
				c.WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000")
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
				},
			},
		},
		{
			name: "Error: Missing relying party UUID",
			before: func(c Client) {
				c.WithRelyingPartyName("DEMO")
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			tt.before(c)

			clientImpl := c.(*client)

			assert.NotNil(t, clientImpl)
			assert.Equal(t, tt.expected.config, clientImpl.config)
		})
	}
}

func Test_WithRelyingPartyName(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "DEMO",
			expected: "DEMO",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithRelyingPartyName(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.RelyingPartyName)
		})
	}
}

func Test_WithRelyingPartyUUID(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "00000000-0000-0000-0000-000000000000",
			expected: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithRelyingPartyUUID(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.RelyingPartyUUID)
		})
	}
}

func Test_WithHashType(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "SHA512",
			expected: "SHA512",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithHashType(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.HashType)
		})
	}
}

func Test_WithText(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "Enter PIN1",
			expected: "Enter PIN1",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithText(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.Text)
		})
	}
}

func TestClient_WithTextFormat(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "GSM-7",
			expected: "GSM-7",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithTextFormat(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.TextFormat)
		})
	}
}

func TestClient_WithLanguage(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "ENG",
			expected: "ENG",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithLanguage(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.Language)
		})
	}
}

func Test_WithURL(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Success",
			param:    "https://tsp.demo.sk.ee/mid-api",
			expected: "https://tsp.demo.sk.ee/mid-api",
		},
		{
			name:     "Empty",
			param:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithURL(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.URL)
		})
	}
}

func Test_WithTimeout(t *testing.T) {
	c := NewClient()

	tests := []struct {
		name     string
		param    time.Duration
		expected time.Duration
	}{
		{
			name:     "Success",
			param:    60 * time.Second,
			expected: 60 * time.Second,
		},
		{
			name:     "Zero",
			param:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = c.WithTimeout(tt.param)
			clientImpl := c.(*client)
			assert.Equal(t, tt.expected, clientImpl.config.Timeout)
		})
	}
}

func TestClient_WithTLSConfig(t *testing.T) {
	manager, err := NewCertificateManager("./certs")
	assert.NoError(t, err)

	tlsConfig := manager.TLSConfig()

	type result struct {
		config *config.Config
	}

	tests := []struct {
		name     string
		before   func(c Client)
		expected result
	}{
		{
			name: "Success",
			before: func(c Client) {
				c.WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
					WithHashType("SHA512").
					WithText("Enter PIN1").
					WithTextFormat("GSM-7").
					WithLanguage("ENG").
					WithURL("https://tsp.demo.sk.ee/mid-api").
					WithTimeout(60 * time.Second).
					WithTLSConfig(tlsConfig)
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
					TLSConfig:        tlsConfig,
				},
			},
		},
		{
			name: "Without TLS Config",
			before: func(c Client) {
				c.
					WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000")
			},
			expected: result{
				config: &config.Config{
					RelyingPartyName: "DEMO",
					RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
					HashType:         "SHA512",
					Text:             "Enter PIN1",
					TextFormat:       "GSM-7",
					Language:         "ENG",
					URL:              "https://tsp.demo.sk.ee/mid-api",
					Timeout:          60 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			tt.before(c)

			clientImpl := c.(*client)
			assert.Equal(t, tt.expected.config, clientImpl.config)
		})
	}
}

func Test_Validate(t *testing.T) {
	tests := []struct {
		name     string
		before   func(c Client)
		expected error
		error    bool
	}{
		{
			name: "Success",
			before: func(c Client) {
				c.
					WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000")
			},
			expected: nil,
			error:    false,
		},
		{
			name: "Error: Missing Relying Party Name",
			before: func(c Client) {
				c.
					WithRelyingPartyName("").
					WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000")
			},
			expected: errors.ErrMissingRelyingPartyName,
			error:    true,
		},
		{
			name: "Error: Missing Relying Party UUID",
			before: func(c Client) {
				c.
					WithRelyingPartyName("DEMO").
					WithRelyingPartyUUID("")
			},
			expected: errors.ErrMissingRelyingPartyUUID,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()

			tt.before(c)

			err := c.Validate()

			if tt.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
