package flibusta

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) Series(ctx context.Context, ID uint64) (*Series, error) {
	uri := "/s/" + strconv.FormatUint(ID, 10)

	content, err := c.readPage(ctx, uri, 0)
	if err != nil {
		return nil, fmt.Errorf("readPage: %w", err)
	}

	booksURLs := rexGroupAll(rexFB2, content, 1)

	for _, page := range totalPages(content) {
		newPageContent, err := c.readPage(ctx, uri, page)
		if err != nil {
			return nil, fmt.Errorf("readPage: %w", err)
		}

		booksURLs = append(booksURLs, rexGroupAll(rexFB2, newPageContent, 1)...)
	}

	return &Series{
		Name:     rexGroup(rexTitle, content, 1),
		BookURLs: booksURLs,
	}, nil
}
