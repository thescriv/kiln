package tezos_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/kiln-mid/pkg/tezos"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock du client HTTP pour simuler les appels
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestTezos_FetchDelegatiosn(t *testing.T) {
	tezosClient := tezos.NewClient()

	for _, tt := range []struct {
		name         string
		mocks        []*gock.Mocker
		queryParams  map[string]any
		tezosOptions tezos.TezosOptions
		response     []tezos.DelegationResponse
		err          error
	}{
		{
			name: "success",
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
			response: []tezos.DelegationResponse{
				{
					ID:     1,
					Amount: 1,
					Level:  1,
					Sender: struct {
						Address string "json:\"address\""
					}{
						Address: "foobar",
					},
					Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC).Format(time.RFC3339),
				},
			},
			tezosOptions: tezos.TezosOptions{},
			err:          nil,
		},
		{
			name: "no delegation found",
			mocks: []*gock.Mocker{
				gock.NewMock(
					gock.NewRequest().URL("https://api.tzkt.io/v1/operations/delegations"),
					gock.NewResponse().BodyString(`[]`).Status(200),
				),
			},
			tezosOptions: tezos.TezosOptions{},
			response:     []tezos.DelegationResponse{},
			err:          nil,
		},
		{
			name:  "no delegation found, invalid From param",
			mocks: []*gock.Mocker{},
			tezosOptions: tezos.TezosOptions{
				From: time.Now().Format(time.April.String()),
			},
			response: []tezos.DelegationResponse{},
			err:      fmt.Errorf("format of field From should be RFC3339"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			gock.DisableNetworking()
			gock.Intercept()

			for _, mock := range tt.mocks {
				gock.Register(mock)
			}

			res, err := tezosClient.FetchDelegations(tt.tezosOptions)
			require.Equal(t, err, tt.err)

			require.Equal(t, res, tt.response)

			require.True(t, gock.IsDone())
		})
	}
}
