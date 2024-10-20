package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/kiln-mid/cmd/worker"
	"github.com/kiln-mid/cmd/xtz"
	"github.com/kiln-mid/pkg/db"
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

	r := gin.Default()

	x := xtz.Handler{
		DelegationsRepository: DelegationsRepository,
	}

	x.RegisterRouter(r)

	tezosClient := tezos.NewClient()

	workerDelegations := worker.NewDelegationWorker(tezosClient, DelegationsRepository)

	utilworker.NewIntervalWorker("worker-delegations", workerDelegations.Run, 0).Start(context.Background())

	r.Run()
}
