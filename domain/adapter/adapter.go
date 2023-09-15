package adapter

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

type API interface {
	GetDelegations(ctx context.Context, startTime time.Time) ([]entity.Delegation, error)
}
