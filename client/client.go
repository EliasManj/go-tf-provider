package client

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HostUrl    string
	HTTPClient *http.Client
}

func NewClient(host string) (*Client, error) {
	_, err := http.Get(host)
	if err != nil {
		return nil, err
	}
	return &Client{
		HostUrl:    host,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
