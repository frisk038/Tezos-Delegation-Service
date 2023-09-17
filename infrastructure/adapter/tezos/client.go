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

	"github.com/frisk038/tezos-delegation-service/domain/entity"
)

// Config represents the configuration for the Tezos API client.
type Config struct {
	Url     string        `yaml:"url" env:"TEZOS-API"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
	Limit   int           `yaml:"limit" env-default:"1"`
}

// httpClient is an interface representing the HTTP client used for making requests.
type httpClient interface {
	Get(url string) (*http.Response, error)
}

// Client is the Tezos API client.
type Client struct {
	Client httpClient
	Url    *url.URL
	Limit  int
}

// delegation is a struct used to parse the response of the Tezos API.
type delegation struct {
	Amount int64  `json:"amount"`
	Block  string `json:"block"`
	Id     int64  `json:"id"`
	Sender struct {
		Address string `json:"address"`
	} `json:"sender"`
	TimeStamp string `json:"timestamp"`
}

// New creates a new instance of the Tezos API client.
func New(cfg Config) (*Client, error) {
	urlApi, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: &http.Client{
			Timeout: cfg.Timeout,
		},
		Url:   urlApi,
		Limit: cfg.Limit,
	}, nil
}

// GetDelegations gets delegations and handles pagination.
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

// getDelegations retrieves delegations from the Tezos API starting from a specified timestamp with an offset to skip what has already been read.
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
		return nil, fmt.Errorf("API returned non-OK status: %v", resp.Status)
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
