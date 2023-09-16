package handler

import (
	"github.com/frisk038/tezos-delegation-service/domain/usecase/delegation"
	"github.com/gin-gonic/gin"
)

func Init(cfg Config, dgUC *delegation.UseCase) *gin.Engine {
	r := gin.Default()
	r.GET("/xtz/delegations", GetDelegations(cfg, dgUC))
	return r
}
