package main

import (
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

	err = cmd.RestCmd(conf)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
