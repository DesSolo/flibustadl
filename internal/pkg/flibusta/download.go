package flibusta

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Download(ctx context.Context, uri string) (*Book, error) {
	resp, err := c.fetch(ctx, uri)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &Book{
		Content:  data,
		FileName: extractFilename(resp),
	}, nil
}

func extractFilename(resp *http.Response) string {
	val := resp.Header.Get("Content-Disposition")

	return strings.Trim(strings.SplitN(val, "=", 2)[1], "\"")
}
