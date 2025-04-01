package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func LoadConfig() (*Config, error) {
	conf := &Config{}
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("no .env file found: %w", err)
	}

	opts := env.Options{
		Prefix:                "APP_",
		UseFieldNameByDefault: true,
	}

	if err := env.ParseWithOptions(conf, opts); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}

	conf.BrokerConfig = loadBrokerConfig()

	return conf, nil
}

func loadBrokerConfig() investgo.Config {
	disableResourceExhaustedRetry, err := strconv.ParseBool(os.Getenv("TINKOFF_DISABLE_RESOURCE_EXHAUSTED_RETRY"))
	if err != nil {
		disableResourceExhaustedRetry = false
		log.Printf("error parsing DISABLE_RESOURCE_EXHAUSTED_RETRY: %f", err)
	}

	disableAllRetry, err := strconv.ParseBool(os.Getenv("TINKOFF_DISABLE_ALL_RETRY"))
	if err != nil {
		disableAllRetry = false
		log.Printf("error parsing TINKOFF_DISABLE_ALL_RETRY: %f", err)
	}

	maxRetries, err := strconv.ParseUint(os.Getenv("TINKOFF_MAX_RETRIES"), 10, 32)
	if err != nil {
		maxRetries = 3
		log.Printf("error parsing TINKOFF_MAX_RETRIES: %f", err)
	}

	return investgo.Config{
		EndPoint:                      os.Getenv("TINKOFF_ADDRESS"),
		Token:                         os.Getenv("TINKOFF_API_KEY"),
		AppName:                       os.Getenv("APP_NAME"),
		AccountId:                     os.Getenv("TINKOFF_ACCOUNT_ID"),
		DisableResourceExhaustedRetry: disableResourceExhaustedRetry,
		DisableAllRetry:               disableAllRetry,
		MaxRetries:                    uint(maxRetries),
	}
}
