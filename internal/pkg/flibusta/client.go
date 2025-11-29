package flibusta

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
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
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, validStatusCode int) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	if resp.StatusCode != validStatusCode {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *Client) fetch(ctx context.Context, uri string) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("newRequest: %w", err)
	}

	return c.do(req, http.StatusOK)
}

func (c *Client) fetchWithRetry(ctx context.Context, uri string) (*http.Response, error) {
	operation := func() (*http.Response, error) {
		return c.fetch(ctx, uri)
	}

	return backoff.Retry(ctx, operation,
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithNotify(func(err error, duration time.Duration) {
			slog.WarnContext(ctx, "backoff.Retry",
				"duration", duration,
				"err", err,
			)
		}),
	)
}

func (c *Client) readPage(ctx context.Context, uri string, page int) ([]byte, error) {
	if page > 0 {
		uri += "?page=" + strconv.Itoa(page)
	}

	slog.DebugContext(ctx, "reading page", slog.String("uri", uri))

	resp, err := c.fetch(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	defer resp.Body.Close()

	return data, nil
}
