package config

import "github.com/russianinvestments/invest-api-go-sdk/investgo"

type Config struct {
	Env          string `env:"ENV" envDefault:".env"`
	AppName      string `env:"APP_NAME" envDefault:"broker-calculator"`
	HTTP         HTTP   `envPrefix:"HTTP_"`
	DB           DB     `envPrefix:"DB_"`
	BrokerConfig investgo.Config
}

type DB struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type HTTP struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port string `env:"PORT" envDefault:"8080"`
}
