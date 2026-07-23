package elasticsearch

import "encoding/json"

type SearchRequest struct {
	Index     string                 `json:"-"`
	Query     map[string]interface{} `json:"query"`
	From      int                    `json:"from,omitempty"`
	Size      int                    `json:"size,omitempty"`
	Sort      []map[string]interface{} `json:"-"`
	Source    interface{}            `json:"_source,omitempty"`
	Highlight map[string]interface{} `json:"highlight,omitempty"`
	RawSort   json.RawMessage        `json:"sort,omitempty"`
}

func (r *SearchRequest) UnmarshalJSON(data []byte) error {
	type Alias SearchRequest
	aux := &struct {
		Sort json.RawMessage `json:"sort,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.Sort) == 0 {
		return nil
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(aux.Sort, &arr); err != nil {
		var obj map[string]interface{}
		if err2 := json.Unmarshal(aux.Sort, &obj); err2 != nil {
			return err
		}
		arr = []map[string]interface{}{obj}
	}
	r.Sort = arr
	return nil
}

type QueryType string

const (
	QueryFunctionScore QueryType = "function_score"
	QueryBool          QueryType = "bool"
	QueryMultiMatch    QueryType = "multi_match"
	QueryMatch         QueryType = "match"
	QueryMatchPhrase   QueryType = "match_phrase"
	QueryTerm          QueryType = "term"
	QueryTerms         QueryType = "terms"
	QueryPrefix        QueryType = "prefix"
	QueryRange         QueryType = "range"
	QueryExists        QueryType = "exists"
)

type BoolQuery struct {
	Must    []map[string]interface{} `json:"must,omitempty"`
	MustNot []map[string]interface{} `json:"must_not,omitempty"`
	Should  []map[string]interface{} `json:"should,omitempty"`
	Filter  []map[string]interface{} `json:"filter,omitempty"`
}

type MultiMatchQuery struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
	Type   string   `json:"type,omitempty"`
	Operator string `json:"operator,omitempty"`
}

type FunctionScoreQuery struct {
	Query      map[string]interface{} `json:"query"`
	FieldValueFactor *FieldValueFactor `json:"field_value_factor,omitempty"`
	BoostMode  string                 `json:"boost_mode,omitempty"`
}

type FieldValueFactor struct {
	Field    string  `json:"field"`
	Factor   float64 `json:"factor,omitempty"`
	Modifier string  `json:"modifier,omitempty"`
	Missing  float64 `json:"missing,omitempty"`
}

type HighlightRequest struct {
	PreTags           []string                          `json:"pre_tags"`
	PostTags          []string                          `json:"post_tags"`
	Fields            map[string]map[string]interface{} `json:"fields"`
	Fragmenter        string                            `json:"fragmenter,omitempty"`
	MaxAnalyzedOffset int                               `json:"max_analyzed_offset,omitempty"`
}
