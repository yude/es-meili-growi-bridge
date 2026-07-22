package meilisearch

type SearchParams struct {
	Q                       string   `json:"q"`
	Offset                  int      `json:"offset,omitempty"`
	Limit                   int      `json:"limit,omitempty"`
	Filter                  string   `json:"filter,omitempty"`
	Sort                    []string `json:"sort,omitempty"`
	AttributesToRetrieve    []string `json:"attributesToRetrieve,omitempty"`
	AttributesToHighlight   []string `json:"attributesToHighlight,omitempty"`
	HighlightPreTag         string   `json:"highlightPreTag,omitempty"`
	HighlightPostTag        string   `json:"highlightPostTag,omitempty"`
	ShowRankingScore        bool     `json:"showRankingScore,omitempty"`
}

type SearchResult struct {
	Hits              []map[string]interface{} `json:"hits"`
	Offset            int                      `json:"offset"`
	Limit             int                      `json:"limit"`
	EstimatedTotalHits int                     `json:"estimatedTotalHits"`
	ProcessingTimeMs  int                      `json:"processingTimeMs"`
	Query             string                   `json:"query"`
}

type IndexInfo struct {
	UID        string `json:"uid"`
	PrimaryKey string `json:"primaryKey"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type IndexStats struct {
	NumberOfDocuments int              `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldDistribution map[string]int  `json:"fieldDistribution"`
}

type Version struct {
	CommitSha string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

type Health struct {
	Status string `json:"status"`
}

type Error struct {
	Message      string `json:"message"`
	Code         string `json:"code"`
	Type         string `json:"type"`
	Link         string `json:"link"`
}

func (e *Error) Error() string {
	return e.Message
}
