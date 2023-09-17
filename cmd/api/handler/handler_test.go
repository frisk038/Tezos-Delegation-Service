package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsecase struct {
	mock.Mock
}

func (mu *mockUsecase) GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error) {
	called := mu.Called(ctx, drq)
	return called.Get(0).([]entity.Delegation), called.Error(1)
}

func getTestContext(method, limit, offset, year string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	u := url.Values{}
	u.Add("limit", limit)
	u.Add("offset", offset)
	u.Add("year", year)
	c.Request.Method = method
	c.Request.URL.RawQuery = u.Encode()

	return c, w
}

func TestGetDelegations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := Config{
		MaxLimit:     100,
		DefaultLimit: 10,
	}
	tn, _ := time.Parse(time.RFC3339, "2023-09-16T11:53:01Z")
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
		c, w := getTestContext("GET", "1", "", "")
		mu := &mockUsecase{}
		mu.On("GetDelegations", c.Request.Context(), entity.DelegationRequest{
			Limit:  1,
			Offset: 0,
			Date:   time.Time{},
		}).Return(dgs, nil)

		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t,
			`{
			"data":[{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1000034,
					"delegator":"dg2",
					"block":"block2"
				},
				{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1234,
					"delegator":"dg1",
					"block":"block1"
				}]
			}`,
			w.Body.String(),
		)
		mu.AssertExpectations(t)
	})

	t.Run("success_with_offset", func(t *testing.T) {
		c, w := getTestContext("GET", "", "1", "")
		mu := &mockUsecase{}
		mu.On("GetDelegations", c.Request.Context(), entity.DelegationRequest{
			Limit:  10,
			Offset: 1,
			Date:   time.Time{},
		}).Return(dgs, nil)

		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t,
			`{
			"data":[{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1000034,
					"delegator":"dg2",
					"block":"block2"
				},
				{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1234,
					"delegator":"dg1",
					"block":"block1"
				}]
			}`,
			w.Body.String(),
		)
		mu.AssertExpectations(t)
	})

	t.Run("success_with_year", func(t *testing.T) {
		c, w := getTestContext("GET", "", "", "2012")
		mu := &mockUsecase{}
		mu.On("GetDelegations", c.Request.Context(), entity.DelegationRequest{
			Limit:  10,
			Offset: 0,
			Date:   time.Date(2012, 01, 01, 0, 0, 0, 0, time.UTC),
		}).Return(dgs, nil)

		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t,
			`{
			"data":[{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1000034,
					"delegator":"dg2",
					"block":"block2"
				},
				{
					"timestamp":"2023-09-16T11:53:01Z",
					"amount":1234,
					"delegator":"dg1",
					"block":"block1"
				}]
			}`,
			w.Body.String(),
		)
		mu.AssertExpectations(t)
	})

	t.Run("wrong_limit_format", func(t *testing.T) {
		c, w := getTestContext("GET", "1s", "", "")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("wrong_offset_format", func(t *testing.T) {
		c, w := getTestContext("GET", "1", "z", "")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("one_digit_year_format", func(t *testing.T) {
		c, w := getTestContext("GET", "1", "", "3")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("letter_year_format", func(t *testing.T) {
		c, w := getTestContext("GET", "1", "", "zzzz")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("wrong_date_year_format", func(t *testing.T) {
		c, w := getTestContext("GET", "1", "", "0000")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("limit_negative", func(t *testing.T) {
		c, w := getTestContext("GET", "-11", "", "")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("limit_over_max", func(t *testing.T) {
		c, w := getTestContext("GET", "1000", "", "")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("offset_negative", func(t *testing.T) {
		c, w := getTestContext("GET", "", "-11", "")
		mu := &mockUsecase{}
		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mu.AssertExpectations(t)
	})

	t.Run("fail_from_uc", func(t *testing.T) {
		c, w := getTestContext("GET", "2", "", "")
		mu := &mockUsecase{}
		mu.On("GetDelegations", c.Request.Context(), entity.DelegationRequest{
			Limit:  2,
			Offset: 0,
			Date:   time.Time{},
		}).Return([]entity.Delegation(nil), errors.New("err"))

		GetDelegations(cfg, mu)(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mu.AssertExpectations(t)
	})
}
