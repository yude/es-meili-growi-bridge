package translator

import (
	"math"

	es "es-meili-growi-bridge/internal/elasticsearch"
	meili "es-meili-growi-bridge/internal/meilisearch"
)

type responseTranslator struct{}

func NewResponseTranslator() ResponseTranslator {
	return &responseTranslator{}
}

func (t *responseTranslator) ToESSearchResponse(meiliResult *meili.SearchResult, index string) *es.SearchResponse {
	if meiliResult == nil {
		return es.NewSearchResponse(0, 0, nil)
	}

	hits := make([]es.Hit, 0, len(meiliResult.Hits))
	for _, hit := range meiliResult.Hits {
		id := extractStringField(hit, "id")
		source := extractSource(hit)
		highlight := extractHighlight(hit)

		score := extractFloat64(hit, "_score")
		if score == 0 {
			score = 1.0
		}
		scorePtr := &score

		hits = append(hits, es.NewHit(index, id, scorePtr, source, highlight))
	}

	total := meiliResult.EstimatedTotalHits
	return es.NewSearchResponse(meiliResult.ProcessingTimeMs, total, hits)
}

func (t *responseTranslator) ToESBulkResponse(indexed, deleted, took int) *es.BulkResponse {
	items := make([]es.BulkItem, 0, indexed+deleted)
	for i := 0; i < indexed; i++ {
		items = append(items, es.NewBulkIndexItem("_index", "_id_placeholder", 201))
	}
	for i := 0; i < deleted; i++ {
		items = append(items, es.NewBulkDeleteItem("_index", "_id_placeholder"))
	}
	return es.NewBulkResponse(took, items)
}

func extractSource(hit map[string]interface{}) map[string]interface{} {
	source := make(map[string]interface{})

	formatted, hasFormatted := hit["_formatted"].(map[string]interface{})

	for k, v := range hit {
		if k == "_formatted" || k == "_matchesPosition" || k == "id" {
			continue
		}
		source[k] = v
	}

	if hasFormatted {
		for k, v := range formatted {
			source[k] = v
		}
	}

	return source
}

func extractHighlight(hit map[string]interface{}) map[string][]string {
	highlight := make(map[string][]string)

	formatted, ok := hit["_formatted"].(map[string]interface{})
	if !ok {
		return nil
	}

	for field, val := range formatted {
		if str, ok := val.(string); ok {
			highlight[field] = []string{str}
			highlight[field+".en"] = []string{str}
			highlight[field+".ja"] = []string{str}
		}
	}

	if len(highlight) == 0 {
		return nil
	}
	return highlight
}

func extractStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func extractFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return 0
}

func float64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	f = math.Round(f*100) / 100
	return &f
}
