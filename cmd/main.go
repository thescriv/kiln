package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/kiln-mid/cmd/xtz"
	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/delegations"
	"github.com/kiln-mid/pkg/tezos"
	"github.com/kiln-mid/pkg/utilconfig"
	"github.com/kiln-mid/pkg/utilworker"
)

func main() {
	utilconfig.LoadConfig()

	dbClient, err := db.CreateClient(os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}

	DelegationsRepository := db.NewDelegationsAdapter(dbClient.DB)

	tezosClient := tezos.NewClient()

	delegationsClient := delegations.NewClient(tezosClient, DelegationsRepository)

	r := gin.Default()

	x := xtz.Handler{
		DelegationsClient: delegationsClient,
	}

	x.RegisterRouter(r)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go utilworker.StartNewIntervalWorker("worker-delegations", func() error {
		delegations, err := delegationsClient.PollNew()
		if err != nil {
			return err
		}

		_, err = delegationsClient.Create(delegations)
		if err != nil {
			return err
		}

		return nil
	}, 0, ctx)

	r.Run()
}
