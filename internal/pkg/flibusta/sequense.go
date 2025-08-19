package flibusta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var (
	rexFB2   = regexp.MustCompile(`<a href="(/b/\d+/fb2)">\(fb2\)</a>`)
	rexTitle = regexp.MustCompile(`<h1 class="title">([^<]+)</h1>`)
)

func (c *Client) Sequence(ctx context.Context, ID uint64) (*Sequence, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/sequence/"+strconv.FormatUint(ID, 10), nil)
	if err != nil {
		return nil, fmt.Errorf("newRequest: %w", err)
	}

	page, err := c.do(req, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	data, err := io.ReadAll(page.Body)
	if err != nil {
		return nil, err
	}
	defer page.Body.Close()

	return &Sequence{
		Name: rexGroup(rexTitle, data, 1),
		URLs: rexGroupAll(rexFB2, data, 1),
	}, nil
}
