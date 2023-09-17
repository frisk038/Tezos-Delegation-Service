package poller

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/adapter"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

type UseCase struct {
	repo repository.Poller
	api  adapter.API
}

func New(repo repository.Poller, api adapter.API) *UseCase {
	return &UseCase{
		repo: repo,
		api:  api,
	}
}

func (uc *UseCase) Fetch(ctx context.Context) error {
	lastDttm, err := uc.repo.SelectLastDelegation(ctx)
	if err != nil {
		return err
	}
	if lastDttm.IsZero() {
		now := time.Now()
		lastDttm = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	dgs, err := uc.api.GetDelegations(ctx, lastDttm)
	if err != nil {
		return err
	}
	if len(dgs) == 0 {
		return nil
	}

	return uc.repo.InsertDelegations(ctx, dgs)
}
