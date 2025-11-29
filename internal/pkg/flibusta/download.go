package flibusta

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) Download(ctx context.Context, uri string) (*Book, error) {
	resp, err := c.fetchWithRetry(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fetchWithRetry: %w", err)
	}

	return &Book{
		Content:  resp.Body,
		FileName: extractFilename(resp),
	}, nil
}

func extractFilename(resp *http.Response) string {
	val := resp.Header.Get("Content-Disposition")

	return strings.Trim(strings.SplitN(val, "=", 2)[1], "\"")
}
