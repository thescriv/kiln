package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/models"
	"github.com/kiln-mid/pkg/tezos"
)

// Worker represent a delegationsWorker
type Worker struct {
	tezosClient           *tezos.Client
	delegationsRepository db.DelegationsRepository
	name                  string
}

// NewDelegationWorker return a new worker handling the delegations part.
// a tezosClient and a DelegationsRepository need to be given in parameter.
func NewDelegationWorker(tezos *tezos.Client, dr db.DelegationsRepository) *Worker {
	return &Worker{
		tezosClient:           tezos,
		delegationsRepository: dr,
		name:                  "worker-delegation",
	}
}

// Run the delegations worker.
// it will poll all recents tezos delegations and insert them in database.
func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	recentDelegations, err := w.delegationsRepository.FindMostRecent(context.Background())
	if err != nil {
		fmt.Printf("delegationsRepository FindMostRecent: %s\n", err)
	}

	var tezosOption = tezos.TezosOptions{
		From: time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
	}

	if recentDelegations != nil {
		tezosOption.From = recentDelegations.Timestamp.Format(time.RFC3339)
	}

	delegationsResponse, err := w.tezosClient.FetchDelegations(tezosOption)
	if err != nil {
		fmt.Printf("tezosClient fetch delegations: %s\n", err)
		return
	}

	if len(delegationsResponse) == 0 {
		return
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

	rowsAffected, err := w.delegationsRepository.CreateMany(context.Background(), &delegations)
	if err != nil {
		fmt.Println("createMany: %w\n", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Printf("[WORKER] %s successfully inserted %d rows.\n", w.name, rowsAffected)
	}
}
