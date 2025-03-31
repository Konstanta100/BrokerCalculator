package main

import (
	"context"
	"github.com/Konstanta100/BrokerCalculator/cmd"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"log"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	ctx := context.Background()
	err = cmd.RestCmd(ctx, conf)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
