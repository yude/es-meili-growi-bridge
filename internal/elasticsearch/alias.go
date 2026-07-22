package elasticsearch

type AliasAction struct {
	Add    *AliasAddRemove `json:"add,omitempty"`
	Remove *AliasAddRemove `json:"remove,omitempty"`
}

type AliasAddRemove struct {
	Alias string `json:"alias"`
	Index string `json:"index"`
}

type UpdateAliasesRequest struct {
	Actions []AliasAction `json:"actions"`
}

type UpdateAliasesResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

func NewUpdateAliasesResponse() *UpdateAliasesResponse {
	return &UpdateAliasesResponse{Acknowledged: true}
}

type PutAliasResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

func NewPutAliasResponse() *PutAliasResponse {
	return &PutAliasResponse{Acknowledged: true}
}
