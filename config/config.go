package config

import (
	"log"
	"os"
	"path"
	"path/filepath"

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
	exec, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exec)

	if err := godotenv.Load(path.Join(dir, ".env")); err != nil {
		log.Fatal(err)
	}
	if err := envconfig.Process("freesound", &Config); err != nil {
		log.Fatal(err)
	}
}
