package translator_test

import (
	"testing"

	es "es-meili-growi-bridge/internal/elasticsearch"
	"es-meili-growi-bridge/internal/translator"
)

func TestSearchTranslator_SimpleQuery(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  "test search",
							"fields": []interface{}{"body", "path"},
							"type":   "most_fields",
						},
					},
				},
			},
		},
		From: 0,
		Size: 20,
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Q != "test search" {
		t.Errorf("expected Q='test search', got '%s'", params.Q)
	}
	if params.Offset != 0 {
		t.Errorf("expected Offset=0, got %d", params.Offset)
	}
	if params.Limit != 20 {
		t.Errorf("expected Limit=20, got %d", params.Limit)
	}
}

func TestSearchTranslator_PhraseQuery(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query": "exact phrase",
							"type":  "phrase",
							"fields": []interface{}{
								"path.raw^2", "body",
							},
						},
					},
				},
			},
		},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Q != `"exact phrase"` {
		t.Errorf("expected Q='\"exact phrase\"', got '%s'", params.Q)
	}
}

func TestSearchTranslator_FilterTerm(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []interface{}{
								map[string]interface{}{
									"term": map[string]interface{}{
										"tag_names": "important",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Filter == "" {
		t.Fatal("expected non-empty filter")
	}
}

func TestSearchTranslator_FunctionScore(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{
								"multi_match": map[string]interface{}{
									"query": "keyword",
									"fields": []interface{}{"body"},
								},
							},
						},
					},
				},
				"field_value_factor": map[string]interface{}{
					"field":    "bookmark_count",
					"modifier": "log1p",
					"factor":   2.5,
					"missing":  0,
				},
				"boost_mode": "sum",
			},
		},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Q != "keyword" {
		t.Errorf("expected Q='keyword', got '%s'", params.Q)
	}
}

func TestSearchTranslator_Sort(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  "test",
							"fields": []interface{}{"body"},
						},
					},
				},
			},
		},
		Sort: []map[string]interface{}{
			{"updated_at": map[string]interface{}{"order": "desc"}},
		},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(params.Sort) != 1 {
		t.Fatalf("expected 1 sort, got %d", len(params.Sort))
	}
	if params.Sort[0] != "updated_at:desc" {
		t.Errorf("expected sort='updated_at:desc', got '%s'", params.Sort[0])
	}
}

func TestSearchTranslator_Highlight(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  "search",
							"fields": []interface{}{"body"},
						},
					},
				},
			},
		},
		Highlight: map[string]interface{}{
			"pre_tags":  []interface{}{"<em>"},
			"post_tags": []interface{}{"</em>"},
			"fields": map[string]interface{}{
				"*": map[string]interface{}{
					"fragment_size": 40,
				},
			},
		},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(params.AttributesToHighlight) == 0 {
		t.Fatal("expected attributesToHighlight to be set")
	}
	if params.HighlightPreTag != "<em>" {
		t.Errorf("expected HighlightPreTag='<em>', got '%s'", params.HighlightPreTag)
	}
	if params.HighlightPostTag != "</em>" {
		t.Errorf("expected HighlightPostTag='</em>', got '%s'", params.HighlightPostTag)
	}
}

func TestSearchTranslator_SourceFiltering(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{
		Query: map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  "test",
							"fields": []interface{}{"body"},
						},
					},
				},
			},
		},
		Source: []interface{}{"path", "body", "bookmark_count"},
	}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(params.AttributesToRetrieve) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(params.AttributesToRetrieve))
	}
}

func TestSearchTranslator_EmptyQuery(t *testing.T) {
	trans := translator.NewSearchTranslator()
	req := &es.SearchRequest{}

	params, err := trans.ToMeilisearchParams(req, "crowi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Q != "" {
		t.Errorf("expected empty Q, got '%s'", params.Q)
	}
}
