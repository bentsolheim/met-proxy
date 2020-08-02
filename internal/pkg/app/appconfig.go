package app

import (
	"github.com/bentsolheim/go-app-utils/utils"
)

type AppConfig struct {
	ServerPort string
}

func ReadAppConfig() AppConfig {
	e := utils.GetEnvOrDefault
	return AppConfig{
		e("SERVER_PORT", "8082"),
	}
}
