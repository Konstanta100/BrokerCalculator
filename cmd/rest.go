package cmd

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/build"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
)

func RestCmd(conf config.Config) error {
	builder := build.New(conf)
	srv, err := builder.RestAPIServer()

	if err != nil {
		return fmt.Errorf("build rest api server: %w", err)
	}

	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
