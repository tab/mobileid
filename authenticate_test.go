package mobileid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tab/mobileid/internal/errors"
)

func Test_CreateSession(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		before      func(w http.ResponseWriter, r *http.Request)
		phoneNumber string
		identity    string
		expected    *Session
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
			expected: &Session{
				Id:   "8fdb516d-1a82-43ba-b82d-be63df569b86",
				Code: "1234",
			},
			err: nil,
		},
		{
			name: "Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected:    nil,
			err:         errors.ErrMobileIdProviderError,
		},
		{
			name: "Internal Server Error",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"title": "Internal Server Error", "status": 500}`))
			},
			phoneNumber: "+37269930366",
			identity:    "51307149560",
			expected:    nil,
			err:         errors.ErrMobileIdProviderError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			client := NewClient()
			client.WithRelyingPartyName("DEMO").
				WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
				WithURL(testServer.URL)

			session, err := client.CreateSession(ctx, tt.phoneNumber, tt.identity)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NotNil(t, session)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_FetchSession(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		before    func(w http.ResponseWriter, r *http.Request)
		sessionId string
		expected  *Person
		err       error
		error     bool
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
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected: &Person{
				IdentityNumber: "PNOEE-51307149560",
				PersonalCode:   "51307149560",
				FirstName:      "MARY ÄNN",
				LastName:       "O'CONNEŽ-ŠUSLIK TESTNUMBER",
			},
			err:   nil,
			error: false,
		},
		{
			name: "Error: Invalid certificate",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`
{
	"state": "COMPLETE",
	"result": "OK",
	"signature": {
		"value": "invalid-signature",
		"algorithm": "sha256WithRSAEncryption"
	},
	"cert": "invalid-certificate",
}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrFailedToDecodeCertificate,
			error:     true,
		},
		{
			name: "Error: Authentication is running",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "RUNNING"}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrAuthenticationIsRunning,
			error:     true,
		},
		{
			name: "Error: USER_CANCELLED",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "USER_CANCELLED"}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       &Error{Code: "USER_CANCELLED"},
			error:     true,
		},
		{
			name: "Error: TIMEOUT",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "TIMEOUT"}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       &Error{Code: "TIMEOUT"},
			error:     true,
		},
		{
			name: "Error: result UNKNOWN",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "COMPLETE", "result": "UNKNOWN"}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrUnsupportedResult,
			error:     true,
		},
		{
			name: "Error: state UNKNOWN",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"state": "UNKNOWN"}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrUnsupportedState,
			error:     true,
		},
		{
			name: "Bad Request",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"title": "Bad Request", "status": 400}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrMobileIdProviderError,
			error:     true,
		},
		{
			name: "Internal Server Error",
			before: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"title": "Internal Server Error", "status": 500}`))
			},
			sessionId: "eb03076a-9f97-423e-af2e-b14c0a481ff9",
			expected:  &Person{},
			err:       errors.ErrMobileIdProviderError,
			error:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.before))
			defer testServer.Close()

			c := NewClient()
			c.WithRelyingPartyName("DEMO").
				WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
				WithURL(testServer.URL)

			session, err := c.FetchSession(ctx, tt.sessionId)

			if tt.error {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NotNil(t, session)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, session)
			}
		})
	}
}

func Test_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Error: USER_CANCELLED",
			code:     "USER_CANCELLED",
			expected: "authentication failed: USER_CANCELLED",
		},
		{
			name:     "Error: SIGNATURE_HASH_MISMATCH",
			code:     "SIGNATURE_HASH_MISMATCH",
			expected: "authentication failed: SIGNATURE_HASH_MISMATCH",
		},
		{
			name:     "Error: PHONE_ABSENT",
			code:     "PHONE_ABSENT",
			expected: "authentication failed: PHONE_ABSENT",
		},
		{
			name:     "Error: DELIVERY_ERROR",
			code:     "DELIVERY_ERROR",
			expected: "authentication failed: DELIVERY_ERROR",
		},
		{
			name:     "Error: SIM_ERROR",
			code:     "SIM_ERROR",
			expected: "authentication failed: SIM_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{Code: tt.code}
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}
