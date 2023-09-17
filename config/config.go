package config

import (
	"github.com/frisk038/tezos-delegation-service/cmd/api/handler"
	"github.com/frisk038/tezos-delegation-service/cmd/cron"
	"github.com/frisk038/tezos-delegation-service/infrastructure/adapter/tezos"
	"github.com/frisk038/tezos-delegation-service/infrastructure/repository"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the configuration structure for the application.
type Config struct {
	Api      handler.Config    `yaml:"api"`
	Tezos    tezos.Config      `yaml:"tezos-client"`
	Database repository.Config `yaml:"database"`
	Cron     cron.Config       `yaml:"cron"`
}

// Cfg is the global configuration instance.
var Cfg Config

// Load reads the configuration from the specified file and populates the global Cfg instance.
// It returns an error if there is an issue reading or parsing the configuration.
func Load(file string) error {
	err := cleanenv.ReadConfig(file, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
