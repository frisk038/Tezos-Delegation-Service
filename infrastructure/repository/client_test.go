package repository

import (
	"context"
	"testing"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestClient_Integration(t *testing.T) {
	tc := initContainer(t)
	c, err := New(Config{ConnUrl: tc.uri}, &slog.Logger{})
	assert.NoError(t, err)
	defer func(){
		c.conn.Close()
		_ = tc.Terminate(context.Background())
	}()

	migrateDb(t, tc)

	for name, fn := range map[string]func(t *testing.T, c *Client){
		"testInsertDelegations" : testInsertDelegations,
	} {
		t.Run(name, func(t *testing.T) {
			fn(t, c)
		})
	}
}

func testInsertDelegations(t *testing.T, c *Client) {
	ctx := context.Background()
	tm := time.Now().UTC()
	dgs := []entity.Delegation{
		{
			Amount: 1000034,
			Block: "block2",
			Id: 3034,
			Delegator: "dg2",
			TimeStamp: tm,
		},
		{
			Amount: 1234,
			Block: "block1",
			Id: 30004,
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
		for rows.Next(){
			var dg entity.Delegation
			err = rows.Scan(&dg.Amount, &dg.Block, &dg.Id, &dg.TimeStamp, &dg.Delegator)
			assert.NoError(t, err)
			got = append(got, dg)
		}
		assert.NoError(t, err)
		assert.Equal(t, dgs, got)
	})
}
