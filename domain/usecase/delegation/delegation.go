package delegation

import (
	"context"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

type UseCase struct {
	repo repository.Delegation
}

func New(repo repository.Delegation) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error) {
	return uc.repo.SelectDelegations(ctx, drq)
}
