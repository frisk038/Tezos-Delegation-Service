package cron

import (
	"context"

	"github.com/robfig/cron/v3"
	"golang.org/x/exp/slog"
)

type Config struct {
	Spec string `yaml:"spec" env-default:"@hourly"`
}

type Cron struct {
	Cr  *cron.Cron
	log *slog.Logger
}

type delegationFetcher interface {
	Fetch(ctx context.Context) error
}

func New(cfg Config, fetcher delegationFetcher, log *slog.Logger) (*Cron, error) {
	c := cron.New()

	_, err := c.AddFunc(cfg.Spec, func() {
		err := fetcher.Fetch(context.Background())
		if err != nil {
			log.Error(err.Error())
		}
	})
	if err != nil {
		return nil, err
	}

	return &Cron{
		Cr:  c,
		log: log,
	}, nil
}
