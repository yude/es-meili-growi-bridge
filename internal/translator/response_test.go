package translator_test

import (
	"testing"

	meili "es-meili-growi-bridge/internal/meilisearch"
	"es-meili-growi-bridge/internal/translator"
)

func TestResponseTranslator_EmptyResult(t *testing.T) {
	trans := translator.NewResponseTranslator()
	result := &meili.SearchResult{
		Hits:              []map[string]interface{}{},
		EstimatedTotalHits: 0,
		ProcessingTimeMs:   3,
	}

	resp := trans.ToESSearchResponse(result, "crowi")

	if resp.Hits.Total.Value != 0 {
		t.Errorf("expected total=0, got %d", resp.Hits.Total.Value)
	}
	if len(resp.Hits.Hits) != 0 {
		t.Errorf("expected 0 hits, got %d", len(resp.Hits.Hits))
	}
	if resp.Took != 3 {
		t.Errorf("expected took=3, got %d", resp.Took)
	}
}

func TestResponseTranslator_WithHits(t *testing.T) {
	trans := translator.NewResponseTranslator()
	result := &meili.SearchResult{
		Hits: []map[string]interface{}{
			{
				"id":    "doc1",
				"title": "Hello",
				"body":  "World",
				"_formatted": map[string]interface{}{
					"body": "<em>World</em>",
				},
			},
		},
		EstimatedTotalHits: 1,
		ProcessingTimeMs:   5,
	}

	resp := trans.ToESSearchResponse(result, "crowi")

	if resp.Hits.Total.Value != 1 {
		t.Errorf("expected total=1, got %d", resp.Hits.Total.Value)
	}
	if len(resp.Hits.Hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(resp.Hits.Hits))
	}

	hit := resp.Hits.Hits[0]
	if hit.ID != "doc1" {
		t.Errorf("expected id='doc1', got '%s'", hit.ID)
	}
	if hit.Index != "crowi" {
		t.Errorf("expected index='crowi', got '%s'", hit.Index)
	}
	if hit.Source["title"] != "Hello" {
		t.Errorf("expected source.title='Hello', got '%v'", hit.Source["title"])
	}

	if hit.Highlight == nil {
		t.Fatal("expected highlight to be set")
	}
	bodyHL, ok := hit.Highlight["body"]
	if !ok {
		t.Fatal("expected body highlight")
	}
	if len(bodyHL) != 1 || bodyHL[0] != "<em>World</em>" {
		t.Errorf("expected body highlight=['<em>World</em>'], got %v", bodyHL)
	}
}

func TestResponseTranslator_NilResult(t *testing.T) {
	trans := translator.NewResponseTranslator()
	resp := trans.ToESSearchResponse(nil, "crowi")

	if resp.Hits.Total.Value != 0 {
		t.Errorf("expected total=0 for nil result, got %d", resp.Hits.Total.Value)
	}
}

func TestResponseTranslator_BulkResponse(t *testing.T) {
	trans := translator.NewResponseTranslator()
	resp := trans.ToESBulkResponse(3, 2, 10)

	if len(resp.Items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(resp.Items))
	}
	if resp.Took != 10 {
		t.Errorf("expected took=10, got %d", resp.Took)
	}

	indexCount := 0
	deleteCount := 0
	for _, item := range resp.Items {
		if item.Index != nil {
			indexCount++
		}
		if item.Delete != nil {
			deleteCount++
		}
	}
	if indexCount != 3 {
		t.Errorf("expected 3 index items, got %d", indexCount)
	}
	if deleteCount != 2 {
		t.Errorf("expected 2 delete items, got %d", deleteCount)
	}
}

func TestResponseTranslator_HighlightMultiFields(t *testing.T) {
	trans := translator.NewResponseTranslator()
	result := &meili.SearchResult{
		Hits: []map[string]interface{}{
			{
				"id":   "doc1",
				"body": "some text",
				"path": "/test",
				"_formatted": map[string]interface{}{
					"body": "<em>some</em> text",
					"path": "<em>/test</em>",
				},
			},
		},
		EstimatedTotalHits: 1,
		ProcessingTimeMs:   2,
	}

	resp := trans.ToESSearchResponse(result, "crowi")

	hit := resp.Hits.Hits[0]
	if hit.Highlight == nil {
		t.Fatal("expected highlight")
	}

	// Should have both body and path entries (and their .en/.ja variants)
	if _, ok := hit.Highlight["body"]; !ok {
		t.Error("expected body highlight key")
	}
	if _, ok := hit.Highlight["path"]; !ok {
		t.Error("expected path highlight key")
	}
}

func TestResponseTranslator_NoFormatted(t *testing.T) {
	trans := translator.NewResponseTranslator()
	result := &meili.SearchResult{
		Hits: []map[string]interface{}{
			{
				"id":    "doc1",
				"title": "Hello",
				"body":  "World",
			},
		},
		EstimatedTotalHits: 1,
		ProcessingTimeMs:   1,
	}

	resp := trans.ToESSearchResponse(result, "crowi")

	hit := resp.Hits.Hits[0]
	if hit.Highlight != nil {
		t.Error("expected nil highlight when no _formatted")
	}
}
