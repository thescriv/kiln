package xtz_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiln-mid/cmd/xtz"
	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/delegations"
	"github.com/kiln-mid/pkg/models"
	"github.com/kiln-mid/pkg/tezos"
	"github.com/kiln-mid/pkg/utilconfig"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
)

func TestGetLastDelegations(t *testing.T) {
	utilconfig.LoadConfig()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	dbClient, err := db.CreateClient(os.Getenv("MYSQL_TEST_DSN"))
	require.NoError(t, err)

	dr := db.NewDelegationsAdapter(dbClient.DB)

	tezosClient := tezos.NewClient()

	delegationsClient := delegations.NewClient(tezosClient, dr)

	dr.CreateMany(&[]models.Delegations{
		{
			ID:        1,
			TezosID:   1,
			Amount:    1,
			Level:     1,
			Delegator: "foobar",
			Timestamp: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			TezosID:   2,
			Amount:    2,
			Level:     2,
			Delegator: "foobar",
			Timestamp: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
		},
	})

	handler := &xtz.Handler{DelegationsClient: delegationsClient}

	handler.RegisterRouter(router)

	tests := []struct {
		name               string
		queryParams        map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:               "Success - No Year",
			queryParams:        map[string]string{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &xtz.Response{Data: []models.Delegations{
				{
					ID:        1,
					TezosID:   1,
					Amount:    1,
					Level:     1,
					Delegator: "foobar",
					Timestamp: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					TezosID:   2,
					Amount:    2,
					Level:     2,
					Delegator: "foobar",
					Timestamp: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
				},
			}, Page: 1},
		},
		{
			name:               "Success - With Year",
			queryParams:        map[string]string{"year": "2023"},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &xtz.Response{Data: []models.Delegations{
				{
					ID:        2,
					TezosID:   2,
					Amount:    2,
					Level:     2,
					Delegator: "foobar",
					Timestamp: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
				},
			}, Page: 1},
		},
		{
			name:               "Error - Invalid Year",
			queryParams:        map[string]string{"year": "1000"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]string{
				"error": "Here are the following available years: 2023,2024",
			},
		},
		{
			name:               "Error - Bad Request",
			queryParams:        map[string]string{"year": "12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"error": "Check Year field is valid and follow the following format `YYYY`",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegationsMarshaled, err := json.Marshal(tt.expectedResponse)
			require.NoError(t, err)

			apitest.New().
				Handler(router).
				Get("/xtz/delegations").
				QueryParams(tt.queryParams).
				Expect(t).
				Status(tt.expectedStatusCode).
				Body(string(delegationsMarshaled)).
				End()
		})
	}
}
