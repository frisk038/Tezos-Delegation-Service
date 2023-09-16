package poller

import (
	"context"

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
		return nil
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
