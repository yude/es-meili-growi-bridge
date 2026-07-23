package translator

import (
	"fmt"
	"strings"

	es "es-meili-growi-bridge/internal/elasticsearch"
	meili "es-meili-growi-bridge/internal/meilisearch"
)

type searchTranslator struct{}

func NewSearchTranslator() SearchTranslator {
	return &searchTranslator{}
}

func (t *searchTranslator) ToMeilisearchParams(req *es.SearchRequest, indexName string) (*meili.SearchParams, error) {
	params := &meili.SearchParams{
		Offset: req.From,
		Limit:  req.Size,
	}

	if req.Source != nil {
		params.AttributesToRetrieve = parseSourceFields(req.Source)
	}

	t.translateSort(req.Sort, params)

	if req.Query != nil {
		t.translateQuery(req.Query, params)
	}

	if req.Highlight != nil {
		t.translateHighlight(req.Highlight, params)
	}

	return params, nil
}

func parseSourceFields(source interface{}) []string {
	switch v := source.(type) {
	case []interface{}:
		fields := make([]string, 0, len(v))
		for _, f := range v {
			if s, ok := f.(string); ok {
				fields = append(fields, s)
			}
		}
		return fields
	case string:
		return []string{v}
	default:
		return nil
	}
}

func (t *searchTranslator) translateSort(sort []map[string]interface{}, params *meili.SearchParams) {
	if len(sort) == 0 {
		return
	}

	for _, s := range sort {
		for field, orderObj := range s {
			if field == "_score" {
				continue
			}
			var order string
			switch v := orderObj.(type) {
			case map[string]interface{}:
				if o, ok := v["order"].(string); ok {
					order = o
				}
			case string:
				order = v
			default:
				order = "desc"
			}
			params.Sort = append(params.Sort, fmt.Sprintf("%s:%s", field, order))
		}
	}
}

func (t *searchTranslator) translateQuery(query map[string]interface{}, params *meili.SearchParams) {
	// Check for function_score wrapper
	q := query
	if fs, ok := query["function_score"].(map[string]interface{}); ok {
		if innerQ, ok := fs["query"].(map[string]interface{}); ok {
			q = innerQ
		}
	}

	boolQ, ok := q["bool"].(map[string]interface{})
	if !ok {
		return
	}

	var qParts []string

	// Build filter expressions
	var filterParts []string

	// Parse must clauses for full-text search terms
	if must, ok := boolQ["must"].([]interface{}); ok {
		qParts = append(qParts, t.parseMustClauses(must)...)
	}

	// Parse must_not clauses for negation
	if mustNot, ok := boolQ["must_not"].([]interface{}); ok {
		negParts := t.parseMustNotClauses(mustNot)
		qParts = append(qParts, negParts...)
	}

	// Parse filter clauses
	if filter, ok := boolQ["filter"].([]interface{}); ok {
		filterParts = t.parseFilterClauses(filter)
	}

	// Build q string
	qStr := strings.TrimSpace(strings.Join(qParts, " "))
	if qStr != "" {
		params.Q = qStr
	}

	// Build filter string
	if len(filterParts) > 0 {
		params.Filter = strings.Join(filterParts, " AND ")
	}
}

func (t *searchTranslator) parseMustClauses(must []interface{}) []string {
	var parts []string
	for _, clause := range must {
		c, ok := clause.(map[string]interface{})
		if !ok {
			continue
		}

		// multi_match query
		if mm, ok := c["multi_match"].(map[string]interface{}); ok {
			queryStr, _ := mm["query"].(string)
			matchType, _ := mm["type"].(string)

			if queryStr == "" {
				continue
			}

			if matchType == "phrase" {
				parts = append(parts, fmt.Sprintf(`"%s"`, queryStr))
			} else {
				parts = append(parts, queryStr)
			}
		}

		// match query
		if match, ok := c["match"].(map[string]interface{}); ok {
			for _, val := range match {
				switch v := val.(type) {
				case string:
					parts = append(parts, v)
				case map[string]interface{}:
					if q, ok := v["query"].(string); ok {
						parts = append(parts, q)
					}
				}
			}
		}

		// match_phrase query
		if mp, ok := c["match_phrase"].(map[string]interface{}); ok {
			for _, val := range mp {
				if q, ok := val.(string); ok {
					parts = append(parts, fmt.Sprintf(`"%s"`, q))
				}
			}
		}
	}
	return parts
}

func (t *searchTranslator) parseMustNotClauses(mustNot []interface{}) []string {
	var parts []string
	for _, clause := range mustNot {
		c, ok := clause.(map[string]interface{})
		if !ok {
			continue
		}

		if mm, ok := c["multi_match"].(map[string]interface{}); ok {
			queryStr, _ := mm["query"].(string)
			matchType, _ := mm["type"].(string)

			if queryStr == "" {
				continue
			}

			if matchType == "phrase" {
				parts = append(parts, fmt.Sprintf(`-"(%s)"`, queryStr))
			} else {
				parts = append(parts, fmt.Sprintf("-%s", queryStr))
			}
		}
	}
	return parts
}

func (t *searchTranslator) parseFilterClauses(filters []interface{}) []string {
	var parts []string
	for _, filter := range filters {
		f, ok := filter.(map[string]interface{})
		if !ok {
			continue
		}

		parts = append(parts, t.parseFilterNode(f)...)
	}
	return parts
}

func (t *searchTranslator) parseFilterNode(node map[string]interface{}) []string {
	// Handle prefix query: { "prefix": { "field": "value" } } or { "prefix": { "field": { "value": "..." } } }
	if prefix, ok := node["prefix"].(map[string]interface{}); ok {
		for field, val := range prefix {
			field = stripFieldSuffix(field)
			valStr := extractValue(val)
			return []string{fmt.Sprintf(`%s STARTS WITH "%s"`, field, escapeFilterValue(valStr))}
		}
	}

	// Handle term query: { "term": { "field": "value" } } or { "term": { "field": { "value": "..." } } }
	if term, ok := node["term"].(map[string]interface{}); ok {
		for field, val := range term {
			if field == "grant" {
				valStr := extractValue(val)
				return []string{fmt.Sprintf("grant = %s", valStr)}
			}
			field = stripFieldSuffix(field)
			valStr := extractValue(val)
			return []string{fmt.Sprintf(`%s = "%s"`, field, escapeFilterValue(valStr))}
		}
	}

	// Handle terms query: { "terms": { "field": ["v1", "v2"] } }
	if terms, ok := node["terms"].(map[string]interface{}); ok {
		for field, val := range terms {
			field = stripFieldSuffix(field)
			values, ok := val.([]interface{})
			if !ok {
				continue
			}
			strVals := make([]string, 0, len(values))
			for _, v := range values {
				strVals = append(strVals, fmt.Sprintf(`"%v"`, v))
			}
			return []string{fmt.Sprintf("%s IN [%s]", field, strings.Join(strVals, ", "))}
		}
	}

	// Handle bool queries (nested)
	if boolQ, ok := node["bool"].(map[string]interface{}); ok {
		var subParts []string

		if must, ok := boolQ["must"].([]interface{}); ok {
			var mustParts []string
			for _, m := range must {
				if mNode, ok := m.(map[string]interface{}); ok {
					mustParts = append(mustParts, t.parseFilterNode(mNode)...)
				}
			}
			if len(mustParts) > 0 {
				if len(mustParts) == 1 {
					subParts = append(subParts, mustParts[0])
				} else {
					subParts = append(subParts, "("+strings.Join(mustParts, " AND ")+")")
				}
			}
		}

		if should, ok := boolQ["should"].([]interface{}); ok {
			var shouldParts []string
			for _, s := range should {
				if sNode, ok := s.(map[string]interface{}); ok {
					shouldParts = append(shouldParts, t.parseFilterNode(sNode)...)
				}
			}
			if len(shouldParts) > 0 {
				if len(shouldParts) == 1 {
					subParts = append(subParts, shouldParts[0])
				} else {
					subParts = append(subParts, "("+strings.Join(shouldParts, " OR ")+")")
				}
			}
		}

		if mustNot, ok := boolQ["must_not"].([]interface{}); ok {
			for _, mn := range mustNot {
				if mnNode, ok := mn.(map[string]interface{}); ok {
					innerParts := t.parseFilterNode(mnNode)
					for _, p := range innerParts {
						subParts = append(subParts, "NOT ("+p+")")
					}
				}
			}
		}

		return subParts
	}

	return nil
}

func (t *searchTranslator) translateHighlight(highlight map[string]interface{}, params *meili.SearchParams) {
	preTags := "<em>"
	postTags := "</em>"

	if pre, ok := highlight["pre_tags"].([]interface{}); ok && len(pre) > 0 {
		if s, ok := pre[0].(string); ok {
			preTags = s
		}
	}
	if post, ok := highlight["post_tags"].([]interface{}); ok && len(post) > 0 {
		if s, ok := post[0].(string); ok {
			postTags = s
		}
	}

	fields := []string{"*"}
	if f, ok := highlight["fields"].(map[string]interface{}); ok {
		fieldNames := make([]string, 0, len(f))
		for name := range f {
			fieldNames = append(fieldNames, name)
		}
		if len(fieldNames) > 0 {
			fields = fieldNames
		}
	}

	params.AttributesToHighlight = fields
	params.HighlightPreTag = preTags
	params.HighlightPostTag = postTags
}

func stripFieldSuffix(field string) string {
	if strings.HasSuffix(field, ".raw") {
		return field[:len(field)-4]
	}
	if strings.HasSuffix(field, ".ja") || strings.HasSuffix(field, ".en") {
		return field[:len(field)-3]
	}
	return field
}

func extractValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case map[string]interface{}:
		if s, ok := val["value"].(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func escapeFilterValue(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
