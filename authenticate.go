package mobileid

import (
	"context"
	"fmt"

	"github.com/tab/mobileid/internal/errors"
	"github.com/tab/mobileid/internal/requests"
	"github.com/tab/mobileid/internal/utils"
)

const (
	Running  = "RUNNING"
	Complete = "COMPLETE"

	OK                      = "OK"
	NOT_MID_CLIENT          = "NOT_MID_CLIENT"
	USER_CANCELLED          = "USER_CANCELLED"
	SIGNATURE_HASH_MISMATCH = "SIGNATURE_HASH_MISMATCH"
	PHONE_ABSENT            = "PHONE_ABSENT"
	DELIVERY_ERROR          = "DELIVERY_ERROR"
	SIM_ERROR               = "SIM_ERROR"
	TIMEOUT                 = "TIMEOUT"
)

// Error represents an error from the Mobile-ID provider
type Error struct {
	Code string
}

// Error returns the error message
func (e *Error) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Code)
}

// CreateSession creates authentication session with the Mobile-ID provider
func (c *client) CreateSession(ctx context.Context, phoneNumber, nationalIdentityNumber string) (*Session, error) {
	session, err := requests.CreateAuthenticationSession(ctx, c.config, phoneNumber, nationalIdentityNumber)
	if err != nil {
		return nil, err
	}

	return (*Session)(session), nil
}

// FetchSession fetches the authentication session from the Mobile-ID provider
func (c *client) FetchSession(ctx context.Context, sessionId string) (*Person, error) {
	response, err := requests.FetchAuthenticationSession(ctx, c.config, sessionId)
	if err != nil {
		return nil, err
	}

	switch response.State {
	case Running:
		return nil, errors.ErrAuthenticationIsRunning
	case Complete:
		switch response.Result {
		case OK:
			person, err := utils.Extract(response.Cert)
			if err != nil {
				return nil, err
			}

			return &Person{
				IdentityNumber: person.IdentityNumber,
				PersonalCode:   person.PersonalCode,
				FirstName:      person.FirstName,
				LastName:       person.LastName,
			}, nil
		case NOT_MID_CLIENT,
			USER_CANCELLED,
			SIGNATURE_HASH_MISMATCH,
			PHONE_ABSENT,
			DELIVERY_ERROR,
			SIM_ERROR,
			TIMEOUT:
			return nil, &Error{Code: response.Result}
		}
	default:
		return nil, errors.ErrUnsupportedState
	}

	return nil, errors.ErrUnsupportedResult
}
