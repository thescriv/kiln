// go:build mage
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/models"
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

	tezosOpt := tezos.TezosOptions{
		From:   time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		To:     time.Date(year, 12, 31, 23, 59, 59, 59, time.UTC).Format(time.RFC3339),
		Limit:  1000,
		Offset: 0,
	}

	for {
		delegationsResponse, err := tezosClient.FetchDelegations(tezosOpt)
		if err != nil {
			panic(err)
		}

		if len(delegationsResponse) == 0 {
			break
		}

		delegations := []models.Delegations{}

		for _, dr := range delegationsResponse {
			if dr.Sender.Address != "" {
				timestamp, err := time.Parse(time.RFC3339, dr.Timestamp)
				if err != nil {
					fmt.Println("time.Parse failed: %w", err)
				}

				d := models.Delegations{
					TezosID:   dr.ID,
					Timestamp: timestamp,
					Level:     dr.Level,
					Amount:    dr.Amount,
					Delegator: dr.Sender.Address,
				}

				delegations = append(delegations, d)
			}
		}

		rowsAffected, err := delegationsRepository.CreateMany(context.Background(), &delegations)
		if err != nil {
			panic(err)
		}

		if rowsAffected > 0 {
			fmt.Printf("[MAGE] FetchDelegationsFromYear successfully inserted %d rows.\n", rowsAffected)
		}

		tezosOpt.Offset += tezosOpt.Limit
	}
}
