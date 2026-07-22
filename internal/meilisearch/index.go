package meilisearch

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetIndex(uid string) (*IndexInfo, error) {
	path := fmt.Sprintf("/indexes/%s", uid)
	resp, err := c.doNoBody("GET", path)
	if err != nil {
		return nil, err
	}

	var info IndexInfo
	if err := jsonUnmarshal(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Client) CreateIndex(uid, primaryKey string) (*IndexInfo, error) {
	body := map[string]string{"uid": uid}
	if primaryKey != "" {
		body["primaryKey"] = primaryKey
	}

	resp, err := c.doRaw("POST", "/indexes", body)
	if err != nil {
		return nil, err
	}

	var info IndexInfo
	if err := jsonUnmarshal(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Client) DeleteIndex(uid string) error {
	path := fmt.Sprintf("/indexes/%s", uid)
	_, err := c.doNoBody("DELETE", path)
	return err
}

func (c *Client) GetIndexStats(uid string) (*IndexStats, error) {
	path := fmt.Sprintf("/indexes/%s/stats", uid)
	resp, err := c.doNoBody("GET", path)
	if err != nil {
		return nil, err
	}

	var stats IndexStats
	if err := jsonUnmarshal(resp, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (c *Client) ListIndexes() ([]IndexInfo, error) {
	resp, err := c.doNoBody("GET", "/indexes")
	if err != nil {
		return nil, err
	}

	var indexes []IndexInfo
	if err := jsonUnmarshal(resp, &indexes); err != nil {
		return nil, err
	}
	return indexes, nil
}

func (c *Client) UpdateFilterableAttributes(uid string, attributes []string) error {
	body := map[string]interface{}{"filterableAttributes": attributes}
	return c.doJSON("PATCH", fmt.Sprintf("/indexes/%s/settings", uid), body, nil)
}

func (c *Client) UpdateSortableAttributes(uid string, attributes []string) error {
	body := map[string]interface{}{"sortableAttributes": attributes}
	return c.doJSON("PATCH", fmt.Sprintf("/indexes/%s/settings", uid), body, nil)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
