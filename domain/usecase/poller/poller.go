package poller

import (
	"context"

	"github.com/frisk038/tezos-delegation-service/domain/adapter"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

type Poller struct {
	repo repository.Poller
	api  adapter.API
}

func New(repo repository.Poller, api adapter.API) *Poller {
	return &Poller{
		repo: repo,
		api:  api,
	}
}

func (p *Poller) Fetch(ctx context.Context) error {
	lastDttm, err := p.repo.SelectLastDelegation(ctx)
	if err != nil {
		return err
	}
	if lastDttm.IsZero() {
		return nil
	}

	dgs, err := p.api.GetDelegations(ctx, lastDttm)
	if err != nil {
		return err
	}
	if len(dgs) == 0 {
		return nil
	}

	return p.repo.InsertDelegations(ctx, dgs)
}
