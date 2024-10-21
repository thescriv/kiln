package delegations

import (
	"fmt"
	"time"

	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/models"
	"github.com/kiln-mid/pkg/tezos"
)

// Client represent the struct of a delegations client.
type Client struct {
	tezosClient           tezos.Client
	delegationsRepository db.DelegationsRepository
}

// NewClient return a new delegations Client to interact with tezos and delegationsRepository.
func NewClient(tezosClient tezos.Client, dr db.DelegationsRepository) *Client {
	return &Client{
		tezosClient:           tezosClient,
		delegationsRepository: dr,
	}
}

// PollWithTezosOptions poll all delegations matching the provided tezosOptions.
func (c Client) PollWithTezosOptions(tezosOpt tezos.TezosDelegationsOption) ([]models.Delegations, error) {
	delegationsResponse, err := c.tezosClient.FetchDelegations(tezosOpt)
	if err != nil {
		return []models.Delegations{}, err
	}

	delegations, err := c.parseDelegations(delegationsResponse)
	if err != nil {
		return []models.Delegations{}, fmt.Errorf("parseDelegations: %w", err)
	}

	return delegations, nil
}

// PollNew poll new delegations based on two possibilities.
// 1. if a delegations is found in database, take the timestamp of the most recent one.
// 2. if no delegation is found in database, start to fetch the earliest delegation from time.Now().AddDate(0, 0, -1)
func (c Client) PollNew() ([]models.Delegations, error) {
	recentDelegations, err := c.delegationsRepository.FindMostRecent()
	if err != nil {
		return []models.Delegations{}, fmt.Errorf("delegationsRepository FindMostRecent: %s", err)
	}

	var tezosOption = tezos.TezosDelegationsOption{
		From: time.Now().AddDate(0, 0, -1),
	}

	if recentDelegations != nil {
		tezosOption.From = recentDelegations.Timestamp
		tezosOption.IDNotIn = append(tezosOption.IDNotIn, recentDelegations.TezosID)
	}

	delegationsResponse, err := c.tezosClient.FetchDelegations(tezosOption)
	if err != nil {
		return []models.Delegations{}, fmt.Errorf("tezosClient fetch delegations: %w", err)
	}

	delegations, err := c.parseDelegations(delegationsResponse)
	if err != nil {
		return []models.Delegations{}, fmt.Errorf("parseDelegations: %w", err)
	}

	return delegations, nil
}

// Create call the delegationsRepository to create given delegations
// number of delegations created are returned.
// return an error if something happen
func (c Client) Create(delegations []models.Delegations) (int64, error) {
	rowsAffected, err := c.delegationsRepository.CreateMany(&delegations)
	if err != nil {
		return rowsAffected, fmt.Errorf("createMany: %w", err)
	}

	return rowsAffected, nil
}

// parseDelegations parse and transform all []tezos.DelegationResponse into a []models.Delegations.
func (c Client) parseDelegations(delegationsResponse []tezos.DelegationResponse) ([]models.Delegations, error) {
	if len(delegationsResponse) == 0 {
		return []models.Delegations{}, nil
	}

	delegations := []models.Delegations{}

	for _, dr := range delegationsResponse {
		if dr.Sender.Address == "" {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, dr.Timestamp)
		if err != nil {
			return []models.Delegations{}, fmt.Errorf("time.Parse failed: %w", err)
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

	return delegations, nil
}
