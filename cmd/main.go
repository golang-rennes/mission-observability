package main

import (
	"context"
	"log"

	"github.com/golang-rennes/mission-observability/api"
	"github.com/golang-rennes/mission-observability/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run

	log.Fatal(api.Run(context.Background(), cfg))
}
