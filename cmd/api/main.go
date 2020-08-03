package main

import (
	"fmt"
	"github.com/bentsolheim/met-proxy/internal/pkg/app"
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()

	cache := app.NewMetCache()
	engine := app.CreateGinEngine(cache)

	return engine.Run(fmt.Sprintf(":%s", config.ServerPort))
}
