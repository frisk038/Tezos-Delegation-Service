package cron

import (
	"context"

	"github.com/robfig/cron/v3"
	"golang.org/x/exp/slog"
)

// Config represents the configuration for the Cron service.
type Config struct {
	Spec string `yaml:"spec" env-default:"@hourly"`
}

// Cron is a service that manages cron jobs.
type Cron struct {
	Cr  *cron.Cron
	log *slog.Logger
}

// delegationFetcher is an interface for fetching delegations.
type delegationFetcher interface {
	Fetch(ctx context.Context) error
}

// New creates a new Cron service with the provided configuration, delegation fetcher, and logger.
// It returns a pointer to the Cron instance and an error if initialization fails.
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
