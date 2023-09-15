package delegation

import (
	"context"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

type DelegationUC struct {
	repo repository.Delegation
}

func New(repo repository.Delegation) *DelegationUC {
	return &DelegationUC{
		repo: repo,
	}
}

func (dg* DelegationUC) GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error){
	return dg.repo.SelectDelegations(ctx, drq)
}
