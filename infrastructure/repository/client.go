package repository

import (
	"context"
	"time"

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
	selectDelegation = `SELECT ts, amount, delegator, block, id
						FROM delegations
						ORDER BY ts DESC
						LIMIT $1
						OFFSET $2;`
	selectLastDelegation = `SELECT ts
							FROM delegations
							ORDER BY ts DESC
							LIMIT 1;`
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

// Select return a slice of delegation from the database, it also handles pagination.
func (c *Client) SelectDelegations(ctx context.Context, limit, offset int) ([]entity.Delegation, error) {
	rows, err := c.conn.Query(ctx, selectDelegation, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entity.Delegation
	for rows.Next() {
		var dg entity.Delegation
		err = rows.Scan(&dg.TimeStamp, &dg.Amount,&dg.Delegator, &dg.Block, &dg.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, dg)
	}

	return res, rows.Err()
}

func (c *Client) SelectLastDelegation(ctx context.Context) (time.Time, error) {
	var lastUpdate time.Time
	err := c.conn.QueryRow(ctx, selectLastDelegation).Scan(&lastUpdate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return lastUpdate, nil
}
