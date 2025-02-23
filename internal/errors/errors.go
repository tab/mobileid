package errors

import (
	"errors"
)

var (
	ErrMissingRelyingPartyName = errors.New("missing required configuration: RelyingPartyName")
	ErrMissingRelyingPartyUUID = errors.New("missing required configuration: RelyingPartyUUID")

	ErrUnsupportedHashType = errors.New("unsupported hash type, allowed hash types are SHA256, SHA384 or SHA512")

	ErrMobileIdProviderError   = errors.New("Mobile-ID provider error")
	ErrMobileIdSessionNotFound = errors.New("Mobile-ID session not found or expired")

	ErrInvalidCertificate    = errors.New("invalid certificate")
	ErrInvalidIdentityNumber = errors.New("invalid identity number")

	ErrFailedToGenerateRandomBytes = errors.New("failed to generate random bytes")

	ErrUnsupportedState  = errors.New("unsupported state, allowed states are COMPLETE or RUNNING")
	ErrUnsupportedResult = errors.New("unsupported result, allowed results are OK or NOT_MID_CLIENT, USER_CANCELLED, SIGNATURE_HASH_MISMATCH, PHONE_ABSENT, DELIVERY_ERROR, SIM_ERROR, TIMEOUT")

	ErrAuthenticationIsRunning = errors.New("authentication is still running")

	ErrFailedToDecodeCertificate = errors.New("failed to decode certificate")
	ErrFailedToParseCertificate  = errors.New("failed to parse certificate")
)
