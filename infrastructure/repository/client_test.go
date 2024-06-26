package repository

import (
	"context"
	"testing"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestClient_Integration(t *testing.T) {
	tc := initContainer(t)
	c, err := New(Config{ConnUrl: tc.uri}, &slog.Logger{})
	assert.NoError(t, err)
	defer func() {
		c.conn.Close()
		_ = tc.Terminate(context.Background())
	}()

	migrateDb(t, tc)

	for name, fn := range map[string]func(t *testing.T, c *Client){
		"testInsertDelegations":    testInsertDelegations,
		"testSelectDelegations":    testSelectDelegations,
		"testSelectLastDelegation": testSelectLastDelegation,
	} {
		t.Run(name, func(t *testing.T) {
			fn(t, c)
			clearTable(context.Background(), t, c.conn)
		})
	}
}

func testInsertDelegations(t *testing.T, c *Client) {
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(time.Millisecond)
	dgs := []entity.Delegation{
		{
			Amount:    1000034,
			Block:     "block2",
			Id:        3034,
			Delegator: "dg2",
			TimeStamp: tm,
		},
		{
			Amount:    1234,
			Block:     "block1",
			Id:        30004,
			Delegator: "dg1",
			TimeStamp: tm,
		},
	}

	t.Run("success", func(t *testing.T) {
		err := c.InsertDelegations(ctx, dgs)
		assert.NoError(t, err)

		rows, err := c.conn.Query(ctx, "SELECT amount, block, id, ts, delegator FROM delegations ORDER by amount DESC")
		assert.NoError(t, err)
		defer rows.Close()
		var got []entity.Delegation
		for rows.Next() {
			var dg entity.Delegation
			err = rows.Scan(&dg.Amount, &dg.Block, &dg.Id, &dg.TimeStamp, &dg.Delegator)
			assert.NoError(t, err)
			got = append(got, dg)
		}
		assert.NoError(t, err)
		assert.Equal(t, dgs, got)
	})
}

func testSelectDelegations(t *testing.T, c *Client) {
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(time.Millisecond)
	dgs := []entity.Delegation{
		{
			Amount:    425,
			Block:     "block4",
			Id:        5600,
			Delegator: "dg4",
			TimeStamp: tm.AddDate(2, 0, 0),
		},
		{
			Amount:    1000034,
			Block:     "block1",
			Id:        3034,
			Delegator: "dg1",
			TimeStamp: tm.Add(2 * time.Minute),
		},
		{
			Amount:    123400,
			Block:     "block2",
			Id:        30004,
			Delegator: "dg2",
			TimeStamp: tm.Add(time.Minute),
		},
		{
			Amount:    100004,
			Block:     "block3",
			Id:        300890,
			Delegator: "dg3",
			TimeStamp: tm,
		},
	}
	require.NoError(t, c.InsertDelegations(ctx, dgs))

	t.Run("success", func(t *testing.T) {
		got, err := c.SelectDelegations(ctx, entity.DelegationRequest{Limit: 5, Offset: 0})
		assert.NoError(t, err)
		assert.Equal(t, dgs, got)
	})

	t.Run("success_with_year", func(t *testing.T) {
		y, err := time.Parse(time.DateOnly, "2023-01-01")
		require.NoError(t, err)
		got, err := c.SelectDelegations(ctx, entity.DelegationRequest{Limit: 5, Offset: 0, Date: y})
		assert.NoError(t, err)
		assert.Equal(t, dgs[1:], got)
	})

	t.Run("success_with_paging", func(t *testing.T) {
		got, err := c.SelectDelegations(ctx, entity.DelegationRequest{Limit: 1, Offset: 0})
		assert.NoError(t, err)
		assert.Equal(t, dgs[:1], got)

		got, err = c.SelectDelegations(ctx, entity.DelegationRequest{Limit: 1, Offset: 1})
		assert.NoError(t, err)
		assert.Equal(t, dgs[1:2], got)
	})

	t.Run("no_rows", func(t *testing.T) {
		clearTable(ctx, t, c.conn)
		got, err := c.SelectDelegations(ctx, entity.DelegationRequest{Limit: 5, Offset: 0})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})
}

func testSelectLastDelegation(t *testing.T, c *Client) {
	ctx := context.Background()
	tm := time.Now().UTC().Truncate(time.Millisecond)
	dgs := []entity.Delegation{
		{
			Amount:    1000034,
			Block:     "block1",
			Id:        3034,
			Delegator: "dg1",
			TimeStamp: tm.Add(time.Minute),
		},
		{
			Amount:    123400,
			Block:     "block2",
			Id:        30004,
			Delegator: "dg2",
			TimeStamp: tm.Add(3 * time.Minute),
		},
		{
			Amount:    100004,
			Block:     "block3",
			Id:        300890,
			Delegator: "dg3",
			TimeStamp: tm,
		},
	}
	require.NoError(t, c.InsertDelegations(ctx, dgs))

	t.Run("success", func(t *testing.T) {
		got, err := c.SelectLastDelegation(ctx)
		assert.NoError(t, err)
		assert.Equal(t, tm.Add(3*time.Minute), got)
	})
	t.Run("no_rows", func(t *testing.T) {
		clearTable(ctx, t, c.conn)
		got, err := c.SelectLastDelegation(ctx)
		assert.NoError(t, err)
		assert.Equal(t, time.Time{}, got)
	})
}
