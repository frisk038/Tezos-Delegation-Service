package delegation

import (
	"context"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/frisk038/tezos-delegation-service/domain/repository"
)

// UseCase represents the use case for managing delegation-related operations.
type UseCase struct {
	repo repository.Delegation // The repository used for delegation data access.
}

// New creates a new instance of the UseCase with the provided delegation repository.
func New(repo repository.Delegation) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// GetDelegations retrieves a list of delegation records based on the specified delegation request.
// It takes a context and a DelegationRequest and returns a slice of delegation entities or an error.
func (uc *UseCase) GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error) {
	return uc.repo.SelectDelegations(ctx, drq)
}
