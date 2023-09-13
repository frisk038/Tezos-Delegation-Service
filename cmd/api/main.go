package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/frisk038/tezos-delegation-service/config"
	"github.com/frisk038/tezos-delegation-service/infrastructure/adapters/tezos"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := config.Load("/Users/olivier/Dev/go/Tezos-Delegation-Service/config/local.yml")
	if err != nil {
		logger.Error(err.Error())
	}

	client, _ := tezos.New()
	testTime0, _ := time.Parse(time.RFC3339, os.Args[1])
	arr, _ := client.GetDelegations(context.Background(), testTime0)
	fmt.Println(arr, len(arr), err)
}
