package auth

import (
	"net/http"

	"github.com/vito-ai/auth/option"
)

type authTransport struct {
	tokenProvider TokenProvider
	transport     http.RoundTripper
}

func NewAuthClient(cliopts *option.ClientOption) (*http.Client, error) {
	tp, err := newAuthTransport(http.DefaultTransport, cliopts)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Transport: tp}
	return httpClient, nil
}

func newAuthTransport(t http.RoundTripper, cliopts *option.ClientOption) (http.RoundTripper, error) {
	tp, err := NewRTZRTokenProvider(cliopts)
	if err != nil {
		return nil, err
	}
	return &authTransport{transport: t, tokenProvider: tp}, nil
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	token, err := t.tokenProvider.Token(req.Context())
	if err != nil {
		return nil, err
	}
	req2 := cloneRequest(req) // per RoundTripper contract
	token.SetAuthHeader(req2)

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

func (t *authTransport) base() http.RoundTripper {
	if t.transport != nil {
		return t.transport
	}
	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
