package credentials

import (
	"os"
)

const (
	ReturnZeroClientId     = "RTZR_CLIENT_ID"
	ReturnZeroClientSecret = "RTZR_CLIENT_SECRET"
)

type ReturnZeroCredentials struct {
	ClientId     string
	ClientSecret string
}

func getVaraiableFromEnv(override string) string {
	return os.Getenv(override)
}

func GetDefaultClientCreds() *ReturnZeroCredentials {
	clientId := getVaraiableFromEnv(ReturnZeroClientId)

	clientSecret := getVaraiableFromEnv(ReturnZeroClientSecret)

	return &ReturnZeroCredentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}
