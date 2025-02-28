package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/tab/mobileid/internal/config"
	"github.com/tab/mobileid/internal/errors"
	"github.com/tab/mobileid/internal/models"
	"github.com/tab/mobileid/internal/utils"
)

const (
	MaxIdleConnections        = 10000
	MaxIdleConnectionsPerHost = 10000
	IdleConnTimeout           = 90 * time.Second
	TLSHandshakeTimeout       = 10 * time.Second
	Timeout                   = 60 * time.Second

	MinMobileIdTimeout = 1000
	MaxMobileIdTimeout = 120000
)

type Response struct {
	Id   string `json:"sessionID"`
	Code string `json:"code"`
}

type Error struct {
	Error   string `json:"error"`
	Time    string `json:"time"`
	TraceId string `json:"traceId"`
}

func CreateAuthenticationSession(
	ctx context.Context,
	cfg *config.Config,
	phoneNumber string,
	identity string,
) (*Response, error) {
	hash, err := utils.GenerateHash(cfg.HashType)
	if err != nil {
		return nil, err
	}

	body := models.AuthenticationRequest{
		RelyingPartyName:       cfg.RelyingPartyName,
		RelyingPartyUUID:       cfg.RelyingPartyUUID,
		PhoneNumber:            phoneNumber,
		NationalIdentityNumber: identity,
		Hash:                   hash,
		HashType:               cfg.HashType,
		Language:               cfg.Language,
		DisplayText:            cfg.Text,
		DisplayTextFormat:      cfg.TextFormat,
	}

	endpoint := fmt.Sprintf("%s/authentication", cfg.URL)
	response, err := httpClient(cfg).R().SetContext(ctx).SetBody(body).Post(endpoint)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		var result Response
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			return nil, err
		}

		code, err := utils.GenerateVerificationCode(hash)
		if err != nil {
			return nil, err
		}

		return &Response{
			Id:   result.Id,
			Code: code,
		}, nil
	case http.StatusBadRequest:
		return nil, errors.ErrMobileIdProviderPayloadError
	case http.StatusUnauthorized:
		return nil, errors.ErrMobileIdAccessForbidden
	case http.StatusMethodNotAllowed:
		return nil, errors.ErrMobileIdMethodNotAllowed
	default:
		return nil, errors.ErrMobileIdProviderError
	}
}

func FetchAuthenticationSession(
	ctx context.Context,
	cfg *config.Config,
	sessionId string,
) (*models.AuthenticationResponse, error) {
	endpoint := fmt.Sprintf("%s/authentication/session/%s", cfg.URL, sessionId)

	timeout := int(cfg.Timeout.Milliseconds())

	switch {
	case timeout < MinMobileIdTimeout:
		timeout = MinMobileIdTimeout
	case timeout > MaxMobileIdTimeout:
		timeout = MaxMobileIdTimeout
	}

	response, err := httpClient(cfg).R().
		SetContext(ctx).
		SetQueryParam("timeoutMs", strconv.Itoa(timeout)).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusOK:
		var result models.AuthenticationResponse
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			return nil, err
		}
		return &result, nil
	case http.StatusForbidden:
		return nil, errors.ErrMobileIdAccessForbidden
	case http.StatusNotFound:
		return nil, errors.ErrMobileIdSessionNotFound
	default:
		return nil, errors.ErrMobileIdProviderError
	}
}

func httpClient(cfg *config.Config) *resty.Client {
	transport := &http.Transport{
		MaxIdleConns:        MaxIdleConnections,
		MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
		IdleConnTimeout:     IdleConnTimeout,
		TLSHandshakeTimeout: TLSHandshakeTimeout,
		TLSClientConfig:     cfg.TLSConfig,
	}

	client := resty.NewWithClient(&http.Client{
		Transport: transport,
		Timeout:   Timeout,
	})

	client.
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	return client
}
