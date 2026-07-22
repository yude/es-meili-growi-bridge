package elasticsearch

type SearchResponse struct {
	Took     int            `json:"took"`
	TimedOut bool           `json:"timed_out"`
	Shards   ShardStatistics `json:"_shards"`
	Hits     SearchHits     `json:"hits"`
	Suggest  interface{}    `json:"suggest,omitempty"`
}

type SearchHits struct {
	Total    SearchTotal `json:"total"`
	MaxScore *float64    `json:"max_score"`
	Hits     []Hit       `json:"hits"`
}

type SearchTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type Hit struct {
	Index     string                 `json:"_index"`
	ID        string                 `json:"_id"`
	Score     *float64               `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Highlight map[string][]string    `json:"highlight,omitempty"`
}

func NewSearchResponse(took int, total int, hits []Hit) *SearchResponse {
	var maxScore *float64
	for _, h := range hits {
		if h.Score != nil {
			if maxScore == nil || *h.Score > *maxScore {
				maxScore = h.Score
			}
		}
	}

	return &SearchResponse{
		Took:     took,
		TimedOut: false,
		Shards:   DefaultShards(),
		Hits: SearchHits{
			Total:    SearchTotal{Value: total, Relation: "eq"},
			MaxScore: maxScore,
			Hits:     hits,
		},
	}
}

func NewHit(index, id string, score *float64, source map[string]interface{}, highlight map[string][]string) Hit {
	return Hit{
		Index:     index,
		ID:        id,
		Score:     score,
		Source:    source,
		Highlight: highlight,
	}
}
