package delegations

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/kiln-mid/pkg/db"
	"github.com/kiln-mid/pkg/miscellaneous"
	"github.com/kiln-mid/pkg/models"
	"github.com/kiln-mid/pkg/tezos"
)

// Client represent the struct of a delegations client.
type Client struct {
	tezosClient           *tezos.Client
	delegationsRepository db.DelegationsRepository
}

// NewClient return a new delegations Client to interact with tezos and delegationsRepository.
func NewClient(tezosClient *tezos.Client, dr db.DelegationsRepository) *Client {
	return &Client{
		tezosClient:           tezosClient,
		delegationsRepository: dr,
	}
}

// GetDelegations return stored delegations based on params received.
// year represent the year to search delegations for, if year is equal to 0 it will retrieve Most Recent delegations, year cannot be equal to something non-present in db.
// page represent the current page for the pagination.
// limit represent the number max of item asked by the client.
func (c Client) GetDelegations(ctx context.Context, year int, page int, limit int) (*[]models.Delegations, error) {
	offset := limit * (page - 1)

	if year == 0 {
		delegations, err := c.delegationsRepository.FindAndOrderByTimestamp(ctx, limit, offset)
		if err != nil {
			return &[]models.Delegations{}, fmt.Errorf("delegationsRepository findAndOrderByTimestamp: %w", err)
		}

		return delegations, nil
	}

	years, err := c.delegationsRepository.FindAvailableYear(ctx)
	if err != nil {
		return &[]models.Delegations{}, fmt.Errorf("delegationsRepository FindAvailableYear: %w", err)
	}

	if !slices.Contains(*years, year) {
		return &[]models.Delegations{}, fmt.Errorf("Here are the following available years: " + miscellaneous.SplitToString(*years, ","))
	}

	delegations, err := c.delegationsRepository.FindFromYear(ctx, year, limit, offset)
	if err != nil {
		return &[]models.Delegations{}, fmt.Errorf("delegationsRepository FindFromYear: %w", err)
	}

	return delegations, nil
}

// PollWithTezosOptions poll all delegations matching the provided tezosOptions.
func (c Client) PollWithTezosOptions(ctx context.Context, tezosOpt tezos.TezosDelegationsOption) ([]models.Delegations, error) {
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
func (c Client) PollNew(ctx context.Context) ([]models.Delegations, error) {
	recentDelegations, err := c.delegationsRepository.FindMostRecent(ctx)
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
func (c Client) Create(ctx context.Context, delegations []models.Delegations) (int64, error) {
	rowsAffected, err := c.delegationsRepository.CreateMany(ctx, &delegations)
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
