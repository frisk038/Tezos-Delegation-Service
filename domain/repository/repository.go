package repository

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

// Poller represents an interface for managing delegation data insertion and retrieval.
type Poller interface {
	InsertDelegations(ctx context.Context, dgs []entity.Delegation) error
	SelectLastDelegation(ctx context.Context) (time.Time, error)
}

// Delegation represents an interface for querying delegation data.
type Delegation interface {
	SelectDelegations(ctx context.Context, dgr entity.DelegationRequest) ([]entity.Delegation, error)
}
