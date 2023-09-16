package main

import (
	"os"

	"github.com/frisk038/tezos-delegation-service/cmd/api/handler"
	"github.com/frisk038/tezos-delegation-service/cmd/cron"
	"github.com/frisk038/tezos-delegation-service/config"
	_ "github.com/frisk038/tezos-delegation-service/docs"
	"github.com/frisk038/tezos-delegation-service/domain/usecase/delegation"
	"github.com/frisk038/tezos-delegation-service/domain/usecase/poller"
	"github.com/frisk038/tezos-delegation-service/infrastructure/adapter/tezos"
	"github.com/frisk038/tezos-delegation-service/infrastructure/repository"
	"golang.org/x/exp/slog"
)

// @title           Tezos Delegation Service
// @version         1.0
// @description     This is a simple service that will poll/return delegations on tezos protocol
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  o.roux2@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  TezosAPI
// @externalDocs.url          https://api.tzkt.io/#operation/Operations_GetDelegations
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := config.Load(os.Args[1])
	if err != nil {
		logger.Error(err.Error())
	}

	db, err := repository.New(config.Cfg.Database, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	tzApi, err := tezos.New()
	if err != nil {
		logger.Error(err.Error())
	}

	pollerUC := poller.New(db, tzApi)
	cr, err := cron.New(config.Cfg.Cron, pollerUC, logger)
	if err != nil {
		logger.Error(err.Error())
	}
	cr.Cr.Start()

	dgUC := delegation.New(db)
	router := handler.Init(config.Cfg.Api, dgUC)
	err = router.Run(config.Cfg.Api.Port)
	if err != nil {
		logger.Error(err.Error())
	}
}
