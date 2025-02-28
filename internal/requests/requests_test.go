package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tab/mobileid/internal/config"
	"github.com/tab/mobileid/internal/errors"
	"github.com/tab/mobileid/internal/models"
)

func Test_CreateAuthenticationSession(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		RelyingPartyName: "DEMO",
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		Text:             "Enter PIN1",
		TextFormat:       "GSM-7",
		Language:         "ENG",
		HashType:         "SHA512",
		Timeout:          10 * time.Second,
	}

	tests := []struct {
		name        string
		before      func(w http.ResponseWriter, r *http.Request)
		phoneNumber string
		identity    string
		expected    *Response
		error       *Error
		err         error
	}{
		{
			name: "Success",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"sessionID": "8fdb516d-1a82-43ba-b82d-be63df569b86", "code": "1234"}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected: &Response{
				Id:   "8fdb516d-1a82-43ba-b82d-be63df569b86",
				Code: "1234",
			},
			err: nil,
		},
		{
			name: "Error: StatusUnauthorized",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"title": "Forbidden", "status": 403}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected:    &Response{},
			err:         errors.ErrMobileIdAccessForbidden,
		},
		{
			name: "Error: Not Found",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"title": "Not Found", "status": 404}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected:    &Response{},
			err:         errors.ErrMobileIdProviderError,
		},
		{
			name: "Error: Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected:    &Response{},
			err:         errors.ErrMobileIdProviderPayloadError,
		},
		{
			name: "Error: Invalid phone number",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
			},
			phoneNumber: "not-a-phone-number",
			identity:    "51307149560",
			err:         errors.ErrMobileIdProviderPayloadError,
			error: &Error{
				Error:   "phoneNumber must contain of + and numbers(8-30)",
				Time:    "2025-02-23T17:31:23",
				TraceId: "d2206fd3aedc3aee",
			},
		},
		{
			name: "Error: Invalid identity number",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
			},
			phoneNumber: "+37269930366",
			identity:    "not-a-personal-code",
			err:         errors.ErrMobileIdProviderPayloadError,
			error: &Error{
				Error:   "nationalIdentityNumber must contain of 11 digits",
				Time:    "2025-02-23T17:40:05",
				TraceId: "65b578c46fb29f6c",
			},
		},
		{
			name: "Error: InternalServerError",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			err:         errors.ErrMobileIdProviderError,
			error: &Error{
				Error:   "Internal Server Error",
				Time:    "2025-02-23T17:40:05",
				TraceId: "65b578c46fb29f6c",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			cfg.URL = testServer.URL

			response, err := CreateAuthenticationSession(ctx, cfg, tt.phoneNumber, tt.identity)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Id, response.Id)
			}
		})
	}
}

func Test_FetchAuthenticationSession(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		RelyingPartyName: "DEMO",
		RelyingPartyUUID: "00000000-0000-0000-0000-000000000000",
		Text:             "Enter PIN1",
		TextFormat:       "GSM-7",
		Language:         "ENG",
		HashType:         "SHA512",
		Timeout:          10 * time.Second,
	}

	id := "8fdb516d-1a82-43ba-b82d-be63df569b86"

	tests := []struct {
		name     string
		before   func(w http.ResponseWriter, r *http.Request)
		id       string
		expected *models.AuthenticationResponse
		err      error
		error    bool
	}{
		{
			name: "Success",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`
{
	"state": "COMPLETE",
	"result": "OK",
	"signature": {
		"value": "Id21arR18nvRSZe3BlJpP1KTwK/wTM3HudXEw3bu/FytpJOrk/i/Lzu+1S47evMFcBON8l4Vw9XNY8M2k9f5yA==",
		"algorithm": "SHA512WithECEncryption"
	},
	"cert": "MIIDqDCCAy6gAwIBAgIQB9W11BzBABj+0d/AZx6UHzAKBggqhkjOPQQDAjBxMQswCQYDVQQGEwJFRTEbMBkGA1UECgwSU0sgSUQgU29sdXRpb25zIEFTMRcwFQYDVQRhDA5OVFJFRS0xMDc0NzAxMzEsMCoGA1UEAwwjVEVTVCBvZiBTSyBJRCBTb2x1dGlvbnMgRUlELVEgMjAyMUUwHhcNMjQwNjEyMDY0NTI4WhcNMjkwNjE2MDY0NTI3WjCBlTELMAkGA1UEBhMCRUUxLzAtBgNVBAMMJk1BUlkgw4ROTixPJ0NPTk5Fxb0txaBVU0xJSyBURVNUTlVNQkVSMSUwIwYDVQQEDBxPJ0NPTk5Fxb0txaBVU0xJSyBURVNUTlVNQkVSMRIwEAYDVQQqDAlNQVJZIMOETk4xGjAYBgNVBAUTEVBOT0VFLTUxMzA3MTQ5NTYwMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWlV1aVSXw6WhagWmFmXE/oe+0R1xZzrHyoiVlgKpGiJ8cwIQLogRGQnWY7NwgQvRHCBmsl99bj57h7SWnd03m6OCAYEwggF9MAkGA1UdEwQCMAAwHwYDVR0jBBgwFoAUScfc7QYUosdtnKbP11L9aOXoBBQwcAYIKwYBBQUHAQEEZDBiMDMGCCsGAQUFBzAChidodHRwOi8vYy5zay5lZS9URVNUX0VJRC1RXzIwMjFFLmRlci5jcnQwKwYIKwYBBQUHMAGGH2h0dHA6Ly9haWEuZGVtby5zay5lZS9laWRxMjAyMWUweAYDVR0gBHEwbzAIBgYEAI96AQIwYwYJKwYBBAHOHxIBMFYwVAYIKwYBBQUHAgEWSGh0dHBzOi8vd3d3LnNraWRzb2x1dGlvbnMuZXUvcmVzb3VyY2VzL2NlcnRpZmljYXRpb24tcHJhY3RpY2Utc3RhdGVtZW50LzA0BgNVHR8ELTArMCmgJ6AlhiNodHRwOi8vYy5zay5lZS90ZXN0X2VpZC1xXzIwMjFlLmNybDAdBgNVHQ4EFgQUj8KjnXvGQJCRYOd5LVfPku7QsZwwDgYDVR0PAQH/BAQDAgeAMAoGCCqGSM49BAMCA2gAMGUCMQCocXWDbBnkM3WEyBdv9Vm0A1MNRv08WrR192dRBcX42Kz5oiH0SdHRJv2ffeuEeSwCMEw2tSA3ClJv233Dl7rIYU/T6UG2NQhvDD5FhnP0umZRmVfAUQ6eVcmU8AhFtNJjwg=="
}`))
			},
			id: id,
			expected: &models.AuthenticationResponse{
				State:  "COMPLETE",
				Result: "OK",
				Signature: models.Signature{
					Value:     "Id21arR18nvRSZe3BlJpP1KTwK/wTM3HudXEw3bu/FytpJOrk/i/Lzu+1S47evMFcBON8l4Vw9XNY8M2k9f5yA==",
					Algorithm: "SHA512WithECEncryption",
				},
				Cert: "MIIDqDCCAy6gAwIBAgIQB9W11BzBABj+0d/AZx6UHzAKBggqhkjOPQQDAjBxMQswCQYDVQQGEwJFRTEbMBkGA1UECgwSU0sgSUQgU29sdXRpb25zIEFTMRcwFQYDVQRhDA5OVFJFRS0xMDc0NzAxMzEsMCoGA1UEAwwjVEVTVCBvZiBTSyBJRCBTb2x1dGlvbnMgRUlELVEgMjAyMUUwHhcNMjQwNjEyMDY0NTI4WhcNMjkwNjE2MDY0NTI3WjCBlTELMAkGA1UEBhMCRUUxLzAtBgNVBAMMJk1BUlkgw4ROTixPJ0NPTk5Fxb0txaBVU0xJSyBURVNUTlVNQkVSMSUwIwYDVQQEDBxPJ0NPTk5Fxb0txaBVU0xJSyBURVNUTlVNQkVSMRIwEAYDVQQqDAlNQVJZIMOETk4xGjAYBgNVBAUTEVBOT0VFLTUxMzA3MTQ5NTYwMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWlV1aVSXw6WhagWmFmXE/oe+0R1xZzrHyoiVlgKpGiJ8cwIQLogRGQnWY7NwgQvRHCBmsl99bj57h7SWnd03m6OCAYEwggF9MAkGA1UdEwQCMAAwHwYDVR0jBBgwFoAUScfc7QYUosdtnKbP11L9aOXoBBQwcAYIKwYBBQUHAQEEZDBiMDMGCCsGAQUFBzAChidodHRwOi8vYy5zay5lZS9URVNUX0VJRC1RXzIwMjFFLmRlci5jcnQwKwYIKwYBBQUHMAGGH2h0dHA6Ly9haWEuZGVtby5zay5lZS9laWRxMjAyMWUweAYDVR0gBHEwbzAIBgYEAI96AQIwYwYJKwYBBAHOHxIBMFYwVAYIKwYBBQUHAgEWSGh0dHBzOi8vd3d3LnNraWRzb2x1dGlvbnMuZXUvcmVzb3VyY2VzL2NlcnRpZmljYXRpb24tcHJhY3RpY2Utc3RhdGVtZW50LzA0BgNVHR8ELTArMCmgJ6AlhiNodHRwOi8vYy5zay5lZS90ZXN0X2VpZC1xXzIwMjFlLmNybDAdBgNVHQ4EFgQUj8KjnXvGQJCRYOd5LVfPku7QsZwwDgYDVR0PAQH/BAQDAgeAMAoGCCqGSM49BAMCA2gAMGUCMQCocXWDbBnkM3WEyBdv9Vm0A1MNRv08WrR192dRBcX42Kz5oiH0SdHRJv2ffeuEeSwCMEw2tSA3ClJv233Dl7rIYU/T6UG2NQhvDD5FhnP0umZRmVfAUQ6eVcmU8AhFtNJjwg==",
			},
			err:   nil,
			error: false,
		},
		{
			name: "Error: USER_CANCELLED",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "USER_CANCELLED"}`))
			},
			id: id,
			expected: &models.AuthenticationResponse{
				State:  "COMPLETE",
				Result: "USER_CANCELLED",
			},
			err:   nil,
			error: false,
		},
		{
			name: "Error: TIMEOUT",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "TIMEOUT"}`))
			},
			id: id,
			expected: &models.AuthenticationResponse{
				State:  "COMPLETE",
				Result: "TIMEOUT",
			},
			err:   nil,
			error: false,
		},
		{
			name: "Error: Not found",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"title": "Not Found", "status": 404}`))
			},
			id:       id,
			expected: &models.AuthenticationResponse{},
			err:      errors.ErrMobileIdSessionNotFound,
			error:    true,
		},
		{
			name: "Error: Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			id:       id,
			expected: &models.AuthenticationResponse{},
			err:      errors.ErrMobileIdProviderError,
			error:    true,
		},
		{
			name: "Error: InternalServerError",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			id:       id,
			expected: &models.AuthenticationResponse{},
			err:      errors.ErrMobileIdProviderError,
			error:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			cfg.URL = testServer.URL

			response, err := FetchAuthenticationSession(ctx, cfg, id)

			if tt.error {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.State, response.State)
				assert.Equal(t, tt.expected.Result, response.Result)
			}
		})
	}
}
