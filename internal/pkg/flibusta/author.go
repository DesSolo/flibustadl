package flibusta

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

func (c *Client) Author(ctx context.Context, ID uint64) (*Author, error) {
	resp, err := c.fetch(ctx, "/a/"+strconv.FormatUint(ID, 10))
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Author{
		Name:     rexGroup(rexTitle, data, 1),
		BookURLs: rexGroupAll(rexFB2, data, 1),
	}, nil
}
