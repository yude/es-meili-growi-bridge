package meilisearch

import (
	"encoding/json"
	"fmt"
)

func (c *Client) AddDocuments(uid string, documents []map[string]interface{}) error {
	path := fmt.Sprintf("/indexes/%s/documents", uid)
	return c.doJSON("POST", path, documents, nil)
}

func (c *Client) DeleteDocuments(uid string, ids []string) error {
	path := fmt.Sprintf("/indexes/%s/documents/delete-batch", uid)
	body := ids
	return c.doJSON("POST", path, body, nil)
}

func (c *Client) DeleteDocument(uid, id string) error {
	path := fmt.Sprintf("/indexes/%s/documents/%s", uid, id)
	_, err := c.doNoBody("DELETE", path)
	return err
}

type BulkResult struct {
	TaskUID int `json:"taskUid"`
}

func (c *Client) AddDocumentsWithPrimaryKey(uid, primaryKey string, documents []map[string]interface{}) error {
	path := fmt.Sprintf("/indexes/%s/documents?primaryKey=%s", uid, primaryKey)
	return c.doJSON("POST", path, documents, nil)
}

func (c *Client) ClearAllDocuments(uid string) error {
	path := fmt.Sprintf("/indexes/%s/documents", uid)
	_, err := c.doNoBody("DELETE", path)
	return err
}

type TaskInfo struct {
	UID       int    `json:"uid"`
	IndexUID  string `json:"indexUid"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	EnqueuedAt string `json:"enqueuedAt"`
}

func (c *Client) WaitForTask(uid string, taskUID int) error {
	path := fmt.Sprintf("/indexes/%s/tasks/%d", uid, taskUID)
	for {
		resp, err := c.doNoBody("GET", path)
		if err != nil {
			return err
		}
		var task TaskInfo
		if err := json.Unmarshal(resp, &task); err != nil {
			return err
		}
		if task.Status == "succeeded" || task.Status == "failed" {
			if task.Status == "failed" {
				return fmt.Errorf("task failed: uid=%d, index=%s", task.UID, task.IndexUID)
			}
			return nil
		}
	}
}
