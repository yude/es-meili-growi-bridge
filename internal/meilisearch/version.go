package meilisearch

func (c *Client) GetVersion() (*Version, error) {
	resp, err := c.doNoBody("GET", "/version")
	if err != nil {
		return nil, err
	}

	var version Version
	if err := jsonUnmarshal(resp, &version); err != nil {
		return nil, err
	}
	return &version, nil
}

func (c *Client) Health() (*Health, error) {
	resp, err := c.doNoBody("GET", "/health")
	if err != nil {
		return nil, err
	}

	var health Health
	if err := jsonUnmarshal(resp, &health); err != nil {
		return nil, err
	}
	return &health, nil
}
