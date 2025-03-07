package cmd

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/Konstanta100/BrokerCalculator/internal/server"
)

func RestCmd(conf config.Config) error {

	srv, err := server.RestAPIServer(conf)

	if err != nil {
		return fmt.Errorf("build rest api server: %w", err)
	}

	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
