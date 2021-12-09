package main

import (
	"encoding/json"
	"fmt"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/met-proxy/internal/pkg/app"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config, err := app.ReadAppConfig()
	if err != nil {
		return stacktrace.Propagate(err, "error while reading application configuration")
	}
	if err := ConfigureLogging("debug"); err != nil {
		return stacktrace.Propagate(err, "unable to configure logging")
	}

	cache := app.NewMetCache()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		requestURI := r.RequestURI
		if data, err := cache.GetFromCacheOrLoad(requestURI); err != nil {
			w.WriteHeader(500)
			err := stacktrace.Propagate(err, "error while getting location data from cache")
			if err := json.NewEncoder(w).Encode(rest.WrapResponse(nil, err)); err != nil {
				log.Error(err)
			}
		} else {
			if _, err := w.Write(data); err != nil {
				log.Error(err)
			}
		}
		log.Debugf("[%d millis] GET %s", time.Since(s)/time.Millisecond, requestURI)
	})
	return http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), nil)
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
