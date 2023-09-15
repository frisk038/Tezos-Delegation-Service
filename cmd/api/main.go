package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/frisk038/tezos-delegation-service/config"
	"github.com/frisk038/tezos-delegation-service/infrastructure/adapter/tezos"
	"github.com/frisk038/tezos-delegation-service/infrastructure/repository"
	"golang.org/x/exp/slog"
)

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

	client, _ := tezos.New()
	testTime0, _ := time.Parse(time.RFC3339, os.Args[2])
	arr, err := client.GetDelegations(context.Background(), testTime0)
	fmt.Println(arr, len(arr), err)

	err = db.InsertDelegations(context.Background(), arr)
	if err != nil {
		logger.Error(err.Error())
	}
}
