package flibusta

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var (
	rexFB2   = regexp.MustCompile(`<a href="(/b/\d+/fb2)">\(fb2\)</a>`)
	rexTitle = regexp.MustCompile(`<h1 class="title">([^<]+)</h1>`)
)

func (c *Client) Sequence(ctx context.Context, ID uint64) (*Sequence, error) {
	resp, err := c.fetch(ctx, "/sequence/"+strconv.FormatUint(ID, 10))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	defer resp.Body.Close()

	return &Sequence{
		Name:     rexGroup(rexTitle, data, 1),
		BookURLs: rexGroupAll(rexFB2, data, 1),
	}, nil
}
