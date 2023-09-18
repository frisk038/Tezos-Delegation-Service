package delegation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (mr *mockRepo) GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error) {
	called := mr.Called(ctx, drq)
	return called.Get(0).([]entity.Delegation), called.Error(1)
}

func (mr *mockRepo) SelectDelegations(ctx context.Context, dgr entity.DelegationRequest) ([]entity.Delegation, error) {
	called := mr.Called(ctx, dgr)
	return called.Get(0).([]entity.Delegation), called.Error(1)
}

func TestUseCase_GetDelegations(t *testing.T) {
	ctx := context.Background()
	tn := time.Now().Truncate(time.Millisecond)
	dgs := []entity.Delegation{
		{
			Amount:    1000034,
			Block:     "block2",
			Id:        3034,
			Delegator: "dg2",
			TimeStamp: tn,
		},
		{
			Amount:    1234,
			Block:     "block1",
			Id:        30004,
			Delegator: "dg1",
			TimeStamp: tn,
		},
	}
	dgr := entity.DelegationRequest{
		Limit: 2,
	}

	t.Run("success", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectDelegations", ctx, dgr).Return(dgs, nil)

		uc := New(mr)
		got, err := uc.GetDelegations(ctx, dgr)

		assert.NoError(t, err)
		assert.Equal(t, dgs, got)
		mr.AssertExpectations(t)
	})
	t.Run("err", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectDelegations", ctx, dgr).Return([]entity.Delegation(nil), errors.New("err"))

		uc := New(mr)
		got, err := uc.GetDelegations(ctx, dgr)

		assert.Error(t, err)
		assert.Nil(t, got)
		mr.AssertExpectations(t)
	})
}
