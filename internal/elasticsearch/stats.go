package elasticsearch

type IndicesStatsResponse struct {
	Shards  ShardStatistics            `json:"_shards"`
	All     map[string]interface{}     `json:"_all"`
	Indices map[string]IndexStats      `json:"indices"`
}

type IndexStats struct {
	UUID      string                  `json:"uuid"`
	Primaries IndexPrimariesStats     `json:"primaries"`
	Total     IndexPrimariesStats     `json:"total"`
}

type IndexPrimariesStats struct {
	Docs     DocsStats     `json:"docs"`
	Store    StoreStats    `json:"store"`
	Indexing IndexingStats `json:"indexing"`
}

type DocsStats struct {
	Count   int `json:"count"`
	Deleted int `json:"deleted"`
}

type StoreStats struct {
	SizeInBytes int64 `json:"size_in_bytes"`
}

type IndexingStats struct {
	IndexTotal           int `json:"index_total"`
	IndexTimeInMillis    int `json:"index_time_in_millis"`
	DeleteTotal          int `json:"delete_total"`
	DeleteTimeInMillis   int `json:"delete_time_in_millis"`
}

func NewEmptyIndexStats() IndexStats {
	return IndexStats{
		UUID: "dummy-uuid",
		Primaries: IndexPrimariesStats{
			Docs:     DocsStats{Count: 0, Deleted: 0},
			Store:    StoreStats{SizeInBytes: 0},
			Indexing: IndexingStats{},
		},
		Total: IndexPrimariesStats{
			Docs:     DocsStats{Count: 0, Deleted: 0},
			Store:    StoreStats{SizeInBytes: 0},
			Indexing: IndexingStats{},
		},
	}
}
