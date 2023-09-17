package handler

import (
	"github.com/frisk038/tezos-delegation-service/domain/usecase/delegation"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Init initializes the Gin HTTP router and sets up the routes.
// It returns a Gin Engine instance that can be used to run the API server.
func Init(cfg Config, dgUC *delegation.UseCase) *gin.Engine {
	r := gin.Default()

	r.GET("/xtz/delegations", GetDelegations(cfg, dgUC))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
