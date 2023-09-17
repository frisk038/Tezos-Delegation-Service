package poller

import (
	"context"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/adapter"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

// UseCase represents the use case for polling and processing delegation data.
type UseCase struct {
	repo repository.Poller // The repository used for delegation data storage.
	api  adapter.API       // The external API adapter for fetching delegation data.
}

// New creates a new instance of the UseCase with the provided repository and API adapter.
func New(repo repository.Poller, api adapter.API) *UseCase {
	return &UseCase{
		repo: repo,
		api:  api,
	}
}

// Fetch retrieves and processes delegation data from an external API.
// It takes a context and returns an error if any operation encounters an error.
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

	// If there are n new delegations, return without further processing.
	if len(dgs) == 0 {
		return nil
	}

	return uc.repo.InsertDelegations(ctx, dgs)
}
