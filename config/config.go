package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type apiConfig struct {
	Port int `yaml:"port" env:"PORT" env-default:"8080"`
}

type tezosConfig struct {
	Url     string        `yaml:"url" env:"TEZOS-API"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
	Limit   int           `yaml:"limit" env-default:"1"`
}

type databaseConfig struct {
	ConnUrl string `yaml:"conn-url" env:"CONNURL"`
}

type cronConfig struct {
	Spec string `yaml:"spec" env-default:"*/10 * * * *"`
}

type Config struct {
	Api      apiConfig      `yaml:"api"`
	Tezos    tezosConfig    `yaml:"tezos-client"`
	Database databaseConfig `yaml:"database"`
	Cron     cronConfig     `yaml:"cron"`
}

var Cfg Config

func Load(file string) error {
	err := cleanenv.ReadConfig(file, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
