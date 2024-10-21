package tezos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/kiln-mid/pkg/miscellaneous"
)

// TezosDelegationsOption represent all option accepted by the tezosClient.
type TezosDelegationsOption struct {
	From    time.Time
	To      time.Time
	IDNotIn []int
	Limit   int
	Offset  int
}

// DelegationResponse represent all value handled by the tezosClient from the endpoint "/v1/operations/delegations".
type DelegationResponse struct {
	ID        int    `json:"id"`
	Level     int    `json:"level"`
	Timestamp string `json:"timestamp"`
	Amount    int    `json:"amount"`
	Sender    struct {
		Address string `json:"address"`
	} `json:"sender"`
}

// createParams return received TezosDelegationsOptions into a url.Values variable.
func (c *Client) createParams(options TezosDelegationsOption) url.Values {
	params := url.Values{}

	if !options.From.IsZero() {
		params.Add("timestamp.ge", options.From.Format(time.RFC3339))
	}

	if !options.To.IsZero() {
		params.Add("timestamp.le", options.To.Format(time.RFC3339))
	}

	if len(options.IDNotIn) > 0 {
		IDs := miscellaneous.SplitToString(options.IDNotIn, ",")
		params.Add("id.ni", IDs)
	}

	if options.Offset != 0 {
		params.Add("offset", strconv.Itoa(options.Offset))
	}

	if options.Limit == 0 {
		options.Limit = 500
	}
	params.Add("limit", strconv.Itoa(options.Limit))

	return params
}

// FetchDelegations fetch all delegations from the endpoint "/v1/operations/delegations"
// Based on the options present in TezosDelegationsOption it will check that `From` and `To` field are in the format of RFC3339
func (c *Client) FetchDelegations(options TezosDelegationsOption) ([]DelegationResponse, error) {
	params := c.createParams(options)

	buffer, err := c.fetcher("/v1/operations/delegations", params)
	if err != nil {
		return []DelegationResponse{}, fmt.Errorf("fetcher: %s", err)
	}

	var d []DelegationResponse

	if len(buffer) > 0 {
		if err := json.NewDecoder(bytes.NewReader(buffer)).Decode(&d); err != nil {
			return nil, fmt.Errorf("body delegations unmarshal: %s", err)
		}
	}

	return d, nil
}
