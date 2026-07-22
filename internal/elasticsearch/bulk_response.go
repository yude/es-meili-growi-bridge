package elasticsearch

type BulkResponse struct {
	Took   int          `json:"took"`
	Errors bool         `json:"errors"`
	Items  []BulkItem   `json:"items"`
}

type BulkItem struct {
	Index  *BulkIndexItem  `json:"index,omitempty"`
	Delete *BulkDeleteItem `json:"delete,omitempty"`
}

type BulkIndexItem struct {
	Index   string `json:"_index"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
}

type BulkDeleteItem struct {
	Index   string `json:"_index"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
}

func NewBulkResponse(took int, items []BulkItem) *BulkResponse {
	hasErrors := false
	for _, item := range items {
		if item.Index != nil && item.Index.Status >= 400 {
			hasErrors = true
		}
		if item.Delete != nil && item.Delete.Status >= 400 {
			hasErrors = true
		}
	}
	return &BulkResponse{
		Took:   took,
		Errors: hasErrors,
		Items:  items,
	}
}

func NewBulkIndexItem(index, id string, status int) BulkItem {
	result := "created"
	if status == 200 {
		result = "updated"
	}
	return BulkItem{
		Index: &BulkIndexItem{
			Index:   index,
			ID:      id,
			Version: 1,
			Result:  result,
			Status:  status,
		},
	}
}

func NewBulkDeleteItem(index, id string) BulkItem {
	return BulkItem{
		Delete: &BulkDeleteItem{
			Index:   index,
			ID:      id,
			Version: 1,
			Result:  "deleted",
			Status:  200,
		},
	}
}
