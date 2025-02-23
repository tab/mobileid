package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	response, err := httpClient().R().SetContext(ctx).SetBody(body).Post(endpoint)
	if err != nil {
		return nil, err
	}

	if response.IsSuccess() {
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
	}

	return nil, errors.ErrMobileIdProviderError
}

func FetchAuthenticationSession(
	ctx context.Context,
	cfg *config.Config,
	sessionId string,
) (*models.AuthenticationResponse, error) {
	endpoint := fmt.Sprintf("%s/authentication/session/%s", cfg.URL, sessionId)

	response, err := httpClient().R().SetContext(ctx).Get(endpoint)
	if err != nil {
		return nil, err
	}

	if response.IsSuccess() {
		var result models.AuthenticationResponse
		if err = json.Unmarshal(response.Body(), &result); err != nil {
			return nil, err
		}

		return &result, nil
	}

	if response.StatusCode() == http.StatusNotFound {
		return nil, errors.ErrMobileIdSessionNotFound
	}

	return nil, errors.ErrMobileIdProviderError
}

func httpClient() *resty.Client {
	transport := &http.Transport{
		MaxIdleConns:        MaxIdleConnections,
		MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
		IdleConnTimeout:     IdleConnTimeout,
		TLSHandshakeTimeout: TLSHandshakeTimeout,
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
