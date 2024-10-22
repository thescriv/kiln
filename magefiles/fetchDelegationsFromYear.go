// go:build mage
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/delegations"
	"github.com/kiln-mid/pkg/tezos"
	"github.com/kiln-mid/pkg/utilconfig"
	"github.com/magefile/mage/mg"
)

type Tezos mg.Namespace

func (Tezos) FetchDelegationsFromYear(ctx context.Context, year int) {
	if year > time.Now().Year() || year < 2018 {
		fmt.Println("Check args : year cannot be before existence of Tezos (2018) and year cannot be in the future")
	}

	utilconfig.LoadConfig()

	dbClient, err := db.CreateClient(os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}

	delegationsRepository := db.NewDelegationsAdapter(dbClient.DB)

	tezosClient := tezos.NewClient()

	delegationClient := delegations.NewClient(tezosClient, delegationsRepository)

	tezosOpt := tezos.TezosDelegationsOption{
		From:   time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		To:     time.Date(year, 12, 31, 23, 59, 59, 59, time.UTC),
		Limit:  1000,
		Offset: 0,
	}

	for {
		delegations, err := delegationClient.PollWithOptions(ctx, tezosOpt)
		if err != nil {
			panic(err)
		}

		if len(delegations) == 0 {
			break
		}

		_, err = delegationClient.Create(ctx, delegations)
		if err != nil {
			panic(err)
		}

		tezosOpt.Offset += tezosOpt.Limit
	}
}
