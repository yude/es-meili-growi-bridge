package translator

import (
	es "es-meili-growi-bridge/internal/elasticsearch"
)

type bulkTranslator struct{}

func NewBulkTranslator() BulkTranslator {
	return &bulkTranslator{}
}

func (t *bulkTranslator) ToMeilisearchOperations(req *es.BulkRequestBody) (indexDocs []map[string]interface{}, deleteIDs []string, indexName string) {
	var indexDocsList []map[string]interface{}
	var deleteIDsList []string
	var commonIndex string

	for _, op := range req.Operations {
		if commonIndex == "" && op.Index != "" {
			commonIndex = op.Index
		}

		switch op.Action {
		case es.BulkActionIndex, es.BulkActionCreate:
			doc := op.Doc
			if doc == nil {
				doc = make(map[string]interface{})
			}
			if op.ID != "" {
				doc["id"] = op.ID
			}
			indexDocsList = append(indexDocsList, doc)

		case es.BulkActionUpdate:
			if op.Doc != nil {
				if op.ID != "" {
					op.Doc["id"] = op.ID
				}
				indexDocsList = append(indexDocsList, op.Doc)
			}

		case es.BulkActionDelete:
			if op.ID != "" {
				deleteIDsList = append(deleteIDsList, op.ID)
			}
		}
	}

	return indexDocsList, deleteIDsList, commonIndex
}
