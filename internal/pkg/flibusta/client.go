package flibusta

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
	url    string
}

func NewClient(url string) *Client {
	return &Client{
		client: http.DefaultClient,
		url:    url,
	}
}

func (c *Client) newRequest(ctx context.Context, method, uri string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.url+uri, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) do(req *http.Request, validStatusCode int) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != validStatusCode {
		return resp, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return resp, nil
}
