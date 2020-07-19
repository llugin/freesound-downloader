package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Config APIConfig

type APIConfig struct {
	ClientID     string `envconfig:"client_id"`
	ClientSecret string `envconfig:"client_secret"`
	ApiKey       string `envconfig:"api_key"`
	Dir          string `envconfig:"dir"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if err := envconfig.Process("freesound", &Config); err != nil {
		log.Fatal(err)
	}
}
