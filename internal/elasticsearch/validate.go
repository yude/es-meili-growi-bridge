package elasticsearch

type ValidateQueryResponse struct {
	Valid        bool              `json:"valid"`
	Shards       ShardStatistics   `json:"_shards"`
	Explanations []Explanation     `json:"explanations,omitempty"`
}

type Explanation struct {
	Index string `json:"index"`
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func NewValidateQueryResponse(valid bool) *ValidateQueryResponse {
	return &ValidateQueryResponse{
		Valid:  valid,
		Shards: DefaultShards(),
	}
}

func NewValidateQueryExplanation(index string, valid bool) *ValidateQueryResponse {
	return &ValidateQueryResponse{
		Valid:  valid,
		Shards: DefaultShards(),
		Explanations: []Explanation{
			{Index: index, Valid: valid},
		},
	}
}
