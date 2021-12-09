package app

import (
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/palantir/stacktrace"
)

type AppConfig struct {
	ServerPort      string
	UserAgentHeader string
}

func ReadAppConfig() (*AppConfig, error) {
	e := utils.GetEnvOrDefault
	s := e("USER_AGENT_HEADER", "")
	if s == "" {
		return nil, stacktrace.NewError("the USER_AGENT_HEADER env variable was not set - you are required to identify youself against api.met.no")
	}
	return &AppConfig{
		e("SERVER_PORT", "8082"),
		s,
	}, nil
}
