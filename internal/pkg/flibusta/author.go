package flibusta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (c *Client) Author(ctx context.Context, ID uint64) (*Author, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/a/"+strconv.FormatUint(ID, 10), nil)
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

	return &Author{
		Name: rexGroup(rexTitle, data, 1),
		URLs: rexGroupAll(rexFB2, data, 1),
	}, nil
}
