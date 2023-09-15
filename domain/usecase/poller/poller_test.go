package poller

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

func (mr *mockRepo) InsertDelegations(ctx context.Context, dgs []entity.Delegation) error {
	return mr.Called(ctx, dgs).Error(0)
}

func (mr *mockRepo) SelectLastDelegation(ctx context.Context) (time.Time, error) {
	called := mr.Called(ctx)
	return called.Get(0).(time.Time), called.Error(1)
}

type mockAPI struct {
	mock.Mock
}

func (ma *mockAPI) GetDelegations(ctx context.Context, startTime time.Time) ([]entity.Delegation, error) {
	called := ma.Called(ctx, startTime)
	return called.Get(0).([]entity.Delegation), called.Error(1)
}

func TestPoller_Fetch(t *testing.T) {
	ctx := context.Background()
	tn := time.Now()
	preTn := tn.Add(-2 * time.Minute)
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

	t.Run("success", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(preTn, nil)
		mr.On("InsertDelegations", ctx, dgs).Return(nil)

		ma := &mockAPI{}
		ma.On("GetDelegations", ctx, preTn).Return(dgs, nil)
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.NoError(t, err)
	})

	t.Run("last_date_empty", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(time.Time{}, nil)

		ma := &mockAPI{}
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.NoError(t, err)
	})

	t.Run("last_date_err", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(time.Time{}, errors.New("err"))

		ma := &mockAPI{}
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.Error(t, err)
	})

	t.Run("api_err", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(preTn, nil)
		mr.On("InsertDelegations", ctx, dgs).Return(nil)

		ma := &mockAPI{}
		ma.On("GetDelegations", ctx, preTn).Return(dgs, errors.New("err"))
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.Error(t, err)
	})

	t.Run("api_empty", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(preTn, nil)

		ma := &mockAPI{}
		ma.On("GetDelegations", ctx, preTn).Return([]entity.Delegation{}, nil)
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.NoError(t, err)
	})

	t.Run("insert_err", func(t *testing.T) {
		mr := &mockRepo{}
		mr.On("SelectLastDelegation", ctx).Return(preTn, nil)
		mr.On("InsertDelegations", ctx, dgs).Return(errors.New("err"))

		ma := &mockAPI{}
		ma.On("GetDelegations", ctx, preTn).Return(dgs, nil)
		p := New(mr, ma)

		err := p.Fetch(ctx)
		assert.Error(t, err)
	})
}
