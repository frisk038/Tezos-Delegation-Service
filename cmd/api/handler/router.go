package handler

import (
	"github.com/frisk038/tezos-delegation-service/domain/usecase/delegation"
	"github.com/gin-gonic/gin"
)

func Init(dgUC *delegation.UseCase) *gin.Engine {
	r := gin.Default()
	r.GET("/xtz/delegations", GetDelegations(dgUC))
	return r
}
