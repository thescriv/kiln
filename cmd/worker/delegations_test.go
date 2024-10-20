package worker_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/kiln-mid/cmd/worker"
	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/models"
	"github.com/kiln-mid/pkg/tezos"
	"github.com/kiln-mid/pkg/utilconfig"
	"github.com/stretchr/testify/require"
)

func TestDelegations_Run(t *testing.T) {
	utilconfig.LoadConfig()

	tezosClient := tezos.NewClient()

	dbClient, err := db.CreateClient(os.Getenv("MYSQL_TEST_DSN"))
	require.NoError(t, err)

	dr := db.NewDelegationsAdapter(dbClient.DB)

	worker := worker.NewDelegationWorker(tezosClient, dr)

	tests := []struct {
		name                string
		insertDelegations   []models.Delegations
		expectedDelegations models.Delegations
		mocks               []*gock.Mocker
	}{
		{
			name:              "find recent delegations (no delegations in base)",
			insertDelegations: []models.Delegations{},
			expectedDelegations: models.Delegations{
				ID:        1,
				TezosID:   1,
				Amount:    1,
				Level:     1,
				Delegator: "foobar",
				Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			mocks: []*gock.Mocker{
				gock.NewMock(
					gock.NewRequest().URL("https://api.tzkt.io/v1/operations/delegations"),
					gock.NewResponse().BodyString(`
						[
							{
								"id": 1,
								"level": 1,
								"timestamp": "2024-01-01T10:00:00Z",
								"sender": {
									"address": "foobar"
								},
								"amount": 1
							}
					]`).Status(200),
				),
			},
		},
		{
			name:                "no recent delegations",
			insertDelegations:   []models.Delegations{},
			expectedDelegations: models.Delegations{},
			mocks: []*gock.Mocker{
				gock.NewMock(
					gock.NewRequest().URL("https://api.tzkt.io/v1/operations/delegations"),
					gock.NewResponse().BodyString(`[]`).Status(200),
				),
			},
		},
		{
			name:              "find recent delegations (insert delegations in base)",
			insertDelegations: []models.Delegations{},
			expectedDelegations: models.Delegations{
				ID:        1,
				TezosID:   1,
				Amount:    1,
				Level:     1,
				Delegator: "foobar",
				Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			mocks: []*gock.Mocker{
				gock.NewMock(
					gock.NewRequest().URL("https://api.tzkt.io/v1/operations/delegations"),
					gock.NewResponse().BodyString(`
						[
							{
								"id": 1,
								"level": 1,
								"timestamp": "2024-01-01T10:00:00Z",
								"sender": {
									"address": "foobar"
								},
								"amount": 1
							}
					]`).Status(200),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbClient.DB.Exec("DELETE FROM delegations")
			dbClient.DB.Exec("TRUNCATE TABLE delegations")

			defer gock.Off()
			gock.DisableNetworking()
			gock.Intercept()

			for _, mock := range tt.mocks {
				gock.Register(mock)
			}

			dr.CreateMany(context.Background(), &tt.insertDelegations)

			var wg sync.WaitGroup

			wg.Add(1)

			worker.Run(&wg)

			wg.Wait()

			d, err := dr.FindMostRecent(context.Background())
			require.NoError(t, err)

			require.Equal(t, d, &tt.expectedDelegations)

			require.True(t, gock.IsDone())
		})
	}
}
