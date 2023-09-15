package repository

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

type Poller interface {
	InsertDelegations(ctx context.Context, dgs []entity.Delegation) error
	SelectLastDelegation(ctx context.Context) (time.Time, error)
}

type Delegation interface {
	SelectDelegations(ctx context.Context, limit, offset int) ([]entity.Delegation, error)
}
