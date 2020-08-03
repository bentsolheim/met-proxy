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
	router := app.CreateGinEngine()

	return router.Run(fmt.Sprintf(":%s", config.ServerPort))
}
