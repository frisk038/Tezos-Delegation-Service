package tezos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/frisk038/tezos-delegation-service/config"
	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

type httpClient interface {
	Get(url string) (*http.Response, error)
}

// Client to the tezos api
type Client struct {
	Client httpClient
	Url    *url.URL
	Limit  int
}

// delegation struct used to parse the response of the api
type delegation struct {
	Amount int64  `json:"amount"`
	Block  string `json:"block"`
	Id     int64  `json:"id"`
	Sender struct {
		Address string `json:"address"`
	} `json:"sender"`
	TimeStamp string `json:"timestamp"`
}

// New creates a new instance of the client
func New() (*Client, error) {
	urlApi, err := url.Parse(config.Cfg.Tezos.Url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: &http.Client{
			Timeout: config.Cfg.Tezos.Timeout,
		},
		Url:   urlApi,
		Limit: config.Cfg.Tezos.Limit,
	}, nil
}

// GetDelegations gets delegation, it also handle pagination
func (c *Client) GetDelegations(ctx context.Context, startTime time.Time) ([]entity.Delegation, error) {
	offset := 0
	var result []entity.Delegation
	for {
		chunk, err := c.getDelegations(ctx, startTime, offset)
		if err != nil {
			return nil, err
		}

		result = append(result, chunk...)
		offset += c.Limit

		if len(chunk) < c.Limit {
			break
		}
	}

	return result, nil
}

// getDelegations get delegations from specified timestamp with an offset to skip what we already read
func (c *Client) getDelegations(_ context.Context, startTime time.Time, offset int) ([]entity.Delegation, error) {
	q := c.Url.Query()
	q.Set("timestamp.gt", startTime.Format(time.RFC3339))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(c.Limit))
	c.Url.RawQuery = q.Encode()

	resp, err := c.Client.Get(c.Url.String())
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned non-OK status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsDgs []delegation
	if err = json.Unmarshal(body, &jsDgs); err != nil {
		return nil, err
	}

	dgs := make([]entity.Delegation, len(jsDgs))
	for i, dg := range jsDgs {
		tm, err := time.Parse(time.RFC3339, dg.TimeStamp)
		if err != nil {
			return nil, err
		}

		dgs[i] = entity.Delegation{
			Amount:    dg.Amount,
			Block:     dg.Block,
			Id:        dg.Id,
			Delegator: dg.Sender.Address,
			TimeStamp: tm,
		}
	}

	return dgs, nil
}
