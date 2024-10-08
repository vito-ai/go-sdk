package option

type ClientOption struct {
	ClientId     string
	ClientSecret string
	Endpoint     string
	TokenURL     string
}

func DefaultClientOption() *ClientOption {
	return &ClientOption{}
}

func (opt *ClientOption) GetRestEndpoint() string {
	if opt.Endpoint != "" {
		return opt.Endpoint
	}
	return "https://openapi.vito.ai/v1/transcribe"
}

func (opt *ClientOption) GetStreamingEndpoint() string {
	if opt.Endpoint != "" {
		return opt.Endpoint
	}
	return "grpc-openapi.vito.ai:443"
}

func (opt *ClientOption) GetTokenURL() string {
	if opt.TokenURL != "" {
		return opt.TokenURL
	}
	return "https://openapi.vito.ai/v1/authenticate"
}

func (opt *ClientOption) GetClientId(override string) string {
	if override != "" {
		return override
	}
	return opt.ClientId
}

func (opt *ClientOption) GetClientSecret(override string) string {
	if override != "" {
		return override
	}
	return opt.ClientSecret
}
