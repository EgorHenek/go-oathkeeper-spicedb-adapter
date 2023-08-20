package configs

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          int    `default:"50150"`
	SpiceDBURL    string `required:"true" envconfig:"SPICE_DB_URL"`
	SpiceDBSecret string `required:"true" envconfig:"SPICE_DB_SECRET"`
}

func NewConfig() Config {
	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
