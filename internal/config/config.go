package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"log"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port string `env:"SERVER_PORT"`
	}
	BrokerConfig investgo.Config
}

func LoadConfig() (Config, error) {
	conf := Config{}

	if err := godotenv.Load(); err != nil {
		return conf, fmt.Errorf("no .env file found: %v", err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		return conf, fmt.Errorf("failed to process env vars: %v", err)
	}

	brokerConfig, err := loadBrokerConfig()

	if err != nil {
		return conf, fmt.Errorf("failed to load broker config: %v", err)
	}

	conf.BrokerConfig = brokerConfig

	return conf, nil
}

func loadBrokerConfig() (investgo.Config, error) {
	var conf investgo.Config

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

	conf = investgo.Config{
		EndPoint:                      os.Getenv("TINKOFF_ADDRESS"),
		Token:                         os.Getenv("TINKOFF_API_KEY"),
		AppName:                       os.Getenv("APP_NAME"),
		AccountId:                     os.Getenv("TINKOFF_ACCOUNT_ID"),
		DisableResourceExhaustedRetry: disableResourceExhaustedRetry,
		DisableAllRetry:               disableAllRetry,
		MaxRetries:                    uint(maxRetries),
	}

	return conf, nil
}

func (c *Config) HTTPAddr() string {
	return net.JoinHostPort(c.Server.Host, c.Server.Port)
}
