package tezos

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockHttp is a mock implementation of an http client
type mockHttp struct {
	mock.Mock
}

func (c *mockHttp) Get(url string) (*http.Response, error) {
	called := c.Called(url)
	return called.Get(0).(*http.Response), called.Error(1)
}

func TestClient_getDelegations(t *testing.T) {
	testTime0, _ := time.Parse(time.RFC3339, "2023-08-01T00:00:00Z")
	testTime1, _ := time.Parse(time.RFC3339, "2023-09-01T00:00:00Z")
	testTime2, _ := time.Parse(time.RFC3339, "2023-09-01T01:00:00Z")
	mockUrl := `=~^https://api\.tzkt\.io/v1/operations/delegations\?limit=.&offset=.&timestamp.gt=2023-08-01T00%3A00%3A00Z`
	apiUrl, _ := url.Parse("https://api.tzkt.io/v1/operations/delegations")

	t.Run("success", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", mockUrl, httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			},
			{
				"amount": 123400,
				"block": "mockBlock2",
				"id": 2,
				"sender": {
					"address": "tz1Sender2"
				},
				"timestamp": "2023-09-01T01:00:00Z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)

		assert.NoError(t, err)
		assert.Len(t, delegations, 2)
		assert.Equal(t, []entity.Delegation{
			{
				Amount:    10023000,
				Block:     "mockBlock1",
				Id:        1,
				Delegator: "tz1Sender1",
				TimeStamp: testTime1,
			},
			{
				Amount:    123400,
				Block:     "mockBlock2",
				Id:        2,
				Delegator: "tz1Sender2",
				TimeStamp: testTime2,
			}},
			delegations)

		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("success_with_paging", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=1&offset=0&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=1&offset=1&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 531000,
				"block": "mockBlock2",
				"id": 2,
				"sender": {
					"address": "tz1Sender2"
				},
				"timestamp": "2023-09-01T01:00:00Z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  1,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)
		assert.NoError(t, err)
		assert.Len(t, delegations, 1)
		assert.Equal(t, []entity.Delegation{
			{
				Amount:    10023000,
				Block:     "mockBlock1",
				Id:        1,
				Delegator: "tz1Sender1",
				TimeStamp: testTime1,
			},
		},
			delegations)

		delegations, err = client.getDelegations(context.Background(), testTime0, 1)
		assert.NoError(t, err)
		assert.Len(t, delegations, 1)
		assert.Equal(t, []entity.Delegation{
			{
				Amount:    531000,
				Block:     "mockBlock2",
				Id:        2,
				Delegator: "tz1Sender2",
				TimeStamp: testTime2,
			},
		},
			delegations)

		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})

	t.Run("success_no_delegations", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", mockUrl, httpmock.NewStringResponder(200, `[]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)

		assert.NoError(t, err)
		assert.Len(t, delegations, 0)
		assert.Equal(t, []entity.Delegation{},
			delegations)

		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("get_err", func(t *testing.T) {
		mh := &mockHttp{}
		mh.On("Get",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=0&timestamp.gt=2023-08-01T00%3A00%3A00Z").
			Return((*http.Response)(nil), errors.New("err"))

		client := &Client{
			Url:    apiUrl,
			Client: mh,
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)
		assert.Error(t, err)
		assert.Nil(t, delegations)
		mh.AssertExpectations(t)
	})

	t.Run("api_not_200", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", mockUrl, httpmock.NewStringResponder(401, `[]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)
		assert.Error(t, err)
		assert.Nil(t, delegations)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("non_RFC333_date", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", mockUrl, httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)

		assert.Error(t, err)
		assert.Nil(t, delegations)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("json_not_valid", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", mockUrl, httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-z"
			},
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.getDelegations(context.Background(), testTime0, 0)

		assert.Error(t, err)
		assert.Nil(t, delegations)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}

func TestClient_GetDelegations(t *testing.T) {
	apiUrl, _ := url.Parse("https://api.tzkt.io/v1/operations/delegations")
	ctx := context.Background()
	testTime0, _ := time.Parse(time.RFC3339, "2023-08-01T00:00:00Z")
	testTime1, _ := time.Parse(time.RFC3339, "2023-09-01T00:00:00Z")

	t.Run("success", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=0&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			},
			{
				"amount": 1093000,
				"block": "mockBlock2",
				"id": 2,
				"sender": {
					"address": "tz1Sender2"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=2&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 531000,
				"block": "mockBlock3",
				"id": 3,
				"sender": {
					"address": "tz1Sender3"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.GetDelegations(ctx, testTime0)
		assert.NoError(t, err)
		assert.Len(t, delegations, 3)
		assert.Equal(t, []entity.Delegation{
			{
				Amount:    10023000,
				Block:     "mockBlock1",
				Id:        1,
				Delegator: "tz1Sender1",
				TimeStamp: testTime1,
			},
			{
				Amount:    1093000,
				Block:     "mockBlock2",
				Id:        2,
				Delegator: "tz1Sender2",
				TimeStamp: testTime1,
			},
			{
				Amount:    531000,
				Block:     "mockBlock3",
				Id:        3,
				Delegator: "tz1Sender3",
				TimeStamp: testTime1,
			},
		}, delegations)

		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})

	t.Run("fail_at_start", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=0&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{,
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			},
			{
				"amount": 1093000,
				"block": "mockBlock2",
				"id": 2,
				"sender": {
					"address": "tz1Sender2"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=2&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 531000,
				"block": "mockBlock3",
				"id": 3,
				"sender": {
					"address": "tz1Sender3"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.GetDelegations(ctx, testTime0)
		assert.Error(t, err)
		assert.Nil(t, delegations)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("fail_in_middle", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=0&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{
				"amount": 10023000,
				"block": "mockBlock1",
				"id": 1,
				"sender": {
					"address": "tz1Sender1"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			},
			{
				"amount": 1093000,
				"block": "mockBlock2",
				"id": 2,
				"sender": {
					"address": "tz1Sender2"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))
		httpmock.RegisterResponder("GET",
			"https://api.tzkt.io/v1/operations/delegations?limit=2&offset=2&timestamp.gt=2023-08-01T00%3A00%3A00Z", httpmock.NewStringResponder(200, `[
			{,
				"amount": 531000,
				"block": "mockBlock3",
				"id": 3,
				"sender": {
					"address": "tz1Sender3"
				},
				"timestamp": "2023-09-01T00:00:00Z"
			}
		]`))

		client := &Client{
			Url:    apiUrl,
			Client: &http.Client{},
			Limit:  2,
		}

		delegations, err := client.GetDelegations(ctx, testTime0)
		assert.Error(t, err)
		assert.Nil(t, delegations)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})
}
