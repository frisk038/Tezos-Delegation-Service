package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/frisk038/tezos-delegation-service/domain/entity"
	"github.com/gin-gonic/gin"
)

type delegationGetter interface {
	GetDelegations(ctx context.Context, drq entity.DelegationRequest) ([]entity.Delegation, error)
}

type Config struct {
	MaxLimit     int `yaml:"max-limit" env:"MAX-LIMIT" env-default:"100"`
	DefaultLimit int `yaml:"default-limit" env:"DEFAULT-LIMIT" env-default:"10"`
}

type delegationJs struct {
	TimeStamp time.Time `json:"timestamp"`
	Amount    int64     `json:"amount"`
	Delegator string    `json:"delegator"`
	Block     string    `json:"block"`
}

// @Summary Get delegations
// @Description Retrieve a list of delegations
// @ID get-delegations
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit the number of results (default is 10)"
// @Param offset query int false "Offset for pagination"
// @Param year query int false "Filter by year (optional)"
// @Success 200 {array} delegationJs
// @Router /xtz/delegations [get]
func GetDelegations(cfg Config, getter delegationGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var tm time.Time
		limit := cfg.DefaultLimit
		offset := 0

		limitRq := c.Query("limit")
		if len(limitRq) != 0 {
			limit, err = strconv.Atoi(limitRq)
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("limit must be numeric"))
				return
			}
			if limit > cfg.MaxLimit || limit < 0 {
				_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("limit must be [0; %d]", cfg.MaxLimit))
				return
			}
		}
		offsetRq := c.Query("offset")
		if len(offsetRq) != 0 {
			offset, err = strconv.Atoi(offsetRq)
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("offset must be numeric"))
				return
			}

			if offset < 0 {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("offset must be positive"))
				return
			}
		}

		yearRq := c.Query("year")
		if len(yearRq) != 0 {
			if len(yearRq) != 4 {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("year must respect XXXX format"))
				return
			}

			year, err := strconv.Atoi(yearRq)
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("year is not a valid number"))
				return
			}

			tm, err = time.Parse(time.DateOnly, fmt.Sprintf("%d-01-01", year))
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("cant format correct date with given year %w", err))
				return
			}
		}

		dgs, err := getter.GetDelegations(c.Request.Context(), entity.DelegationRequest{
			Limit:  limit,
			Offset: offset,
			Date:   tm,
		})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		resp := []delegationJs{}
		for _, dg := range dgs {
			resp = append(resp, delegationJs{
				TimeStamp: dg.TimeStamp,
				Amount:    dg.Amount,
				Delegator: dg.Delegator,
				Block:     dg.Block,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": resp})
	}
}
