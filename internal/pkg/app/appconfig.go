package app

import (
	"github.com/bentsolheim/go-app-utils/utils"
	"strconv"
)

type AppConfig struct {
	ServerPort          string
	SkipTlsVerification bool
}

func ReadAppConfig() AppConfig {
	e := utils.GetEnvOrDefault
	skipTlsVerification, err := strconv.ParseBool(e("SKIP_TLS_VERIFICATION", "false"))
	if err != nil {
		skipTlsVerification = false
	}
	return AppConfig{
		e("SERVER_PORT", "8082"),
		skipTlsVerification,
	}
}
