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

type dgJs struct {
	TimeStamp time.Time `json:"timestamp"`
	Amount    int64     `json:"amount"`
	Delegator string    `json:"delegator"`
	Block     string    `json:"block"`
}

func GetDelegations(getter delegationGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var tm time.Time
		limit := 10
		offset := 0

		limitRq := c.Query("limit")
		if len(limitRq) != 0 {
			limit, err = strconv.Atoi(limitRq)
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("limit must be numeric"))
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

			//2023-09-15T15:00:00
			tm, err = time.Parse(time.RFC3339, fmt.Sprintf("%d-01-01", year))
			if err != nil {
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("cant format correct date with given year"))
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

		var resp []dgJs
		for _, dg := range dgs {
			resp = append(resp, dgJs{
				TimeStamp: dg.TimeStamp,
				Amount:    dg.Amount,
				Delegator: dg.Delegator,
				Block:     dg.Block,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": resp})
	}
}
