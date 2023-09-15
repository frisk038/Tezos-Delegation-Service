package repository

import (
	"context"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

// Config represents the configuration for the database connection.
type Config struct {
	ConnUrl string `yaml:"conn-url" env:"CONNURL"`
}

// Client is a PostgreSQL database client for handling delegations.
type Client struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

const (
	insertDelegation = `INSERT INTO delegations
						(id, ts, amount, delegator, block)
						VALUES ($1, $2, $3, $4, $5);`
)

// New creates a new PostgreSQL client for handling delegations.
func New(cfg Config, logger *slog.Logger) (*Client, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.ConnUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: dbpool,
		log:  logger,
	}, nil
}

// InsertDelegations inserts a batch of delegations into the database.
func (c *Client) InsertDelegations(ctx context.Context, dgs []entity.Delegation) error {
	batch := &pgx.Batch{}
	for _, dg := range dgs {
		// Queue each delegation for insertion using the prepared SQL statement.
		batch.Queue(insertDelegation, dg.Id, dg.TimeStamp, dg.Amount, dg.Delegator, dg.Block)
	}

	br := c.conn.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return err
	}

	return nil
}
