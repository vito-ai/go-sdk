package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vito-ai/auth/credentials"
	"github.com/vito-ai/auth/option"
)

// TokenProvider Interface
type TokenProvider interface {
	Token(context.Context) (*ReturnZeroToken, error)
}

// Default Token Provider for RTZR
type tokenProviderRTZR struct {
	token        *ReturnZeroToken
	Client       *http.Client
	clientId     string
	clientSecret string
	TokenURL     string
}

func NewRTZRTokenProvider(opt *option.ClientOption) (TokenProvider, error) {
	creds := credentials.GetDefaultClientCreds()
	tp := &tokenProviderRTZR{
		clientId:     opt.GetClientId(creds.ClientId),
		clientSecret: opt.GetClientSecret(creds.ClientSecret),
		TokenURL:     opt.GetTokenURL(),
		Client:       http.DefaultClient,
	}

	if err := tp.validate(); err != nil {
		return nil, err
	}

	return tp, nil
}

func (o *tokenProviderRTZR) validate() error {
	if o == nil {
		return errors.New("auth : options must be provided")
	}
	if o.clientId == "" {
		return errors.New("auth: RTZR_CLIENT_ID must be provided")
	}
	if o.clientSecret == "" {
		return errors.New("auth: RTZR_CLIENT_SECRET must be provided")
	}
	if o.TokenURL == "" {
		return errors.New("auth: TokenURL must be provided")
	}
	return nil
}

func (tp *tokenProviderRTZR) Token(ctx context.Context) (*ReturnZeroToken, error) {
	if tp.token.isValidWithExpiry() {
		return tp.token, nil
	}

	formData := url.Values{}
	formData.Set("client_id", tp.clientId)
	formData.Set("client_secret", tp.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tp.TokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tp.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when making authentication request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from authentication server: %s", resp.Status)
	}

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	result := &ReturnZeroToken{}
	if err = json.Unmarshal(respByte, result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	if result.AccessToken == "" {
		return nil, errors.New("unmarshalled response is missing access_token")
	}
	if result.ExpireAt <= time.Now().Unix() {
		return nil, errors.New("unmarshalled response has invalid expire_at timestamp")
	}

	tp.token = result
	return result, nil
}
