package cron

import (
	"github.com/robfig/cron/v3"
	"golang.org/x/exp/slog"
)

type Config struct {
	Spec string `yaml:"spec" env-default:"@hourly"`
}

type Cron struct {
	cr  *cron.Cron
	log *slog.Logger
}

func New(cfg Config, log *slog.Logger) (*Cron, error) {
	c := cron.New()
	_, err := c.AddFunc(cfg.Spec, func() {
		println("runing...")
	})
	if err != nil {
		return nil, err
	}

	return &Cron{
		cr:  c,
		log: log,
	}, nil
}
