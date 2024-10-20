package tezos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kiln-mid/pkg/miscellaneous"
	"github.com/kiln-mid/pkg/utilhttp"
)

// Client represent a tezosClient.
type Client struct {
	HTTP    utilhttp.Client
	BaseUrl string
}

// NewClient create a new http client to handle request on `https://api.tzkt.io/` api.
func NewClient() *Client {
	client := &Client{
		HTTP:    utilhttp.NewClient(45 * time.Second),
		BaseUrl: `https://api.tzkt.io/`,
	}

	return client
}

// TezosOptions represent all option accepted by the tezosClient.
type TezosOptions struct {
	From    string
	To      string
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

// fetcher is the tezosFetcher, it is a private function as only tezosClient should call this function.
func (c *Client) fetcher(path string, params url.Values) ([]byte, error) {
	u, _ := url.ParseRequestURI(c.BaseUrl)
	u.Path = path
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http Do: %s", err)
	}
	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io ReadAll: %s", err)
	}

	return buffer, err
}

// FetchDelegations fetch all delegations from the endpoint "/v1/operations/delegations"
// Based on the options present in TezosOptions it will check that `From` and `To` field are in the format of RFC3339
func (c *Client) FetchDelegations(options TezosOptions) ([]DelegationResponse, error) {
	params := url.Values{}

	if options.From != "" {
		_, err := time.Parse(time.RFC3339, options.From)
		if err != nil {
			return []DelegationResponse{}, fmt.Errorf("format of field From should be RFC3339")
		}

		params.Add("timestamp.ge", options.From)
	}

	if options.To != "" {
		_, err := time.Parse(time.RFC3339, options.To)
		if err != nil {
			return []DelegationResponse{}, fmt.Errorf("format of field To should be RFC3339")
		}

		params.Add("timestamp.le", options.To)
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
