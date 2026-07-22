package translator

import (
	es "es-meili-growi-bridge/internal/elasticsearch"
	meili "es-meili-growi-bridge/internal/meilisearch"
)

type SearchTranslator interface {
	ToMeilisearchParams(req *es.SearchRequest, indexName string) (*meili.SearchParams, error)
}

type BulkTranslator interface {
	ToMeilisearchOperations(req *es.BulkRequestBody) (indexDocs []map[string]interface{}, deleteIDs []string, indexName string)
}

type ResponseTranslator interface {
	ToESSearchResponse(meiliResult *meili.SearchResult, index string) *es.SearchResponse
	ToESBulkResponse(indexed int, deleted int, took int) *es.BulkResponse
}
