package meilisearch

import (
	"fmt"
)

func (c *Client) Search(uid string, params *SearchParams) (*SearchResult, error) {
	path := fmt.Sprintf("/indexes/%s/search", uid)

	resp, err := c.doRaw("POST", path, params)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := jsonUnmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("decode search result: %w", err)
	}
	return &result, nil
}
