package tezos

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

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
