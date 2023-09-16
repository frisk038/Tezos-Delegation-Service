package config

import (
	"time"

	"github.com/frisk038/tezos-delegation-service/cmd/api/handler"
	"github.com/frisk038/tezos-delegation-service/cmd/cron"
	"github.com/frisk038/tezos-delegation-service/infrastructure/repository"
	"github.com/ilyakaznacheev/cleanenv"
)

type tezosConfig struct {
	Url     string        `yaml:"url" env:"TEZOS-API"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
	Limit   int           `yaml:"limit" env-default:"1"`
}

type Config struct {
	Api      handler.Config    `yaml:"api"`
	Tezos    tezosConfig       `yaml:"tezos-client"`
	Database repository.Config `yaml:"database"`
	Cron     cron.Config       `yaml:"cron"`
}

var Cfg Config

func Load(file string) error {
	err := cleanenv.ReadConfig(file, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
