package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateHash(t *testing.T) {

	tests := []struct {
		name     string
		hashType string
		error    bool
	}{
		{
			name:     "Success: SHA256",
			hashType: "SHA256",
			error:    false,
		},
		{
			name:     "Success: SHA384",
			hashType: "SHA384",
			error:    false,
		},
		{
			name:     "Success: SHA512",
			hashType: "SHA512",
			error:    false,
		},
		{
			name:     "Error: Unsupported hash type",
			hashType: "MD5",
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GenerateHash(tt.hashType)

			if tt.error {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func Test_VerificationCode(t *testing.T) {
	tests := []struct {
		name  string
		hash  string
		error bool
	}{
		{
			name:  "Success",
			hash:  "aGVsbG8=",
			error: false,
		},
		{
			name:  "Error: Invalid base64",
			hash:  "aGVsbG8",
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := GenerateVerificationCode(tt.hash)

			if tt.error {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
			}
		})
	}
}
