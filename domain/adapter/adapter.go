package adapter

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

// API is an interface that defines the methods for interacting with the Tezos API.
type API interface {
	GetDelegations(ctx context.Context, startTime time.Time) ([]entity.Delegation, error)
}
