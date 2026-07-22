package elasticsearch

type IndexCreateRequest struct {
	Index    string                 `json:"-"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Mappings map[string]interface{} `json:"mappings,omitempty"`
	Aliases  map[string]interface{} `json:"aliases,omitempty"`
}

type IndexCreateResponse struct {
	Index       string `json:"index"`
	ShardsAck   bool   `json:"shards_acknowledged"`
	Acknowledged bool  `json:"acknowledged"`
}

func NewIndexCreateResponse(index string) *IndexCreateResponse {
	return &IndexCreateResponse{
		Index:        index,
		ShardsAck:    true,
		Acknowledged: true,
	}
}

type IndexDeleteResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

func NewIndexDeleteResponse() *IndexDeleteResponse {
	return &IndexDeleteResponse{Acknowledged: true}
}

type IndexExistsResponse bool

type IndexAliasInfo struct {
	Aliases map[string]interface{} `json:"aliases"`
}
