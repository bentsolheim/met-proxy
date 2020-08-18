package main

import (
	"fmt"
	"github.com/bentsolheim/met-proxy/internal/pkg/app"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()
	if err := ConfigureLogging("debug"); err != nil {
		return err
	}

	cache := app.NewMetCache(config.SkipTlsVerification)
	engine := app.CreateGinEngine(cache)

	return engine.Run(fmt.Sprintf(":%s", config.ServerPort))
}

func ConfigureLogging(logLevel string) error {
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return stacktrace.Propagate(err, "error while parsing log level %s", logLevel)
	}
	log.SetLevel(level)
	return nil
}
