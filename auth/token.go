package auth

import (
	"fmt"
	"net/http"
	"time"
)

type ReturnZeroToken struct {
	AccessToken string `json:"access_token"`
	ExpireAt    int64  `json:"expire_at"`
}

func (t *ReturnZeroToken) isEmpty() bool {
	return t == nil || t.AccessToken == ""
}

func (t *ReturnZeroToken) isValidWithExpiry() bool {
	if t.isEmpty() {
		return false
	}
	if t.ExpireAt < time.Now().Unix() {
		return false
	}
	return true
}

func (t ReturnZeroToken) SetAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("%s %v", "Bearer", t.AccessToken))
}
