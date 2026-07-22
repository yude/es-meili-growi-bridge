package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"fmt"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) Bulk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to read request body"))
		return
	}

	operations, err := parseNDJSON(body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to parse NDJSON: "+err.Error()))
		return
	}

	batchReq := &elasticsearch.BulkRequestBody{Operations: operations}

	indexDocs, deleteIDs, indexName := h.bulkTrans.ToMeilisearchOperations(batchReq)

	indexedCount := 0
	deletedCount := 0

	if len(indexDocs) > 0 {
		meiliIndex := resolveIndex(h, indexName)
		if meiliIndex == "" {
			meiliIndex = indexName
		}
		if err := h.meiliClient.AddDocuments(meiliIndex, indexDocs); err != nil {
			writeError(w, elasticsearch.NewBadRequest("bulk index error: "+err.Error()))
			return
		}
		indexedCount = len(indexDocs)
	}

	if len(deleteIDs) > 0 {
		meiliIndex := resolveIndex(h, indexName)
		if meiliIndex == "" {
			meiliIndex = indexName
		}
		if err := h.meiliClient.DeleteDocuments(meiliIndex, deleteIDs); err != nil {
			writeError(w, elasticsearch.NewBadRequest("bulk delete error: "+err.Error()))
			return
		}
		deletedCount = len(deleteIDs)
	}

	bulkResp := h.responseTrans.ToESBulkResponse(indexedCount, deletedCount, 0)
	writeJSON(w, http.StatusOK, bulkResp)
}

type bulkAction struct {
	Index *struct {
		Index string `json:"_index"`
		ID    string `json:"_id"`
		Type  string `json:"_type,omitempty"`
	} `json:"index,omitempty"`
	Delete *struct {
		Index string `json:"_index"`
		ID    string `json:"_id"`
		Type  string `json:"_type,omitempty"`
	} `json:"delete,omitempty"`
}

func parseNDJSON(data []byte) ([]elasticsearch.BulkOperation, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var operations []elasticsearch.BulkOperation
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lineNum++

		var action bulkAction
		if err := json.Unmarshal([]byte(line), &action); err != nil {
			return nil, err
		}

		var op elasticsearch.BulkOperation

		switch {
		case action.Index != nil:
			op.Action = elasticsearch.BulkActionIndex
			op.Index = action.Index.Index
			op.ID = action.Index.ID

			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("unexpected end of NDJSON after index action")
			}

			docLine := strings.TrimSpace(scanner.Text())
			if docLine == "" {
				return nil, fmt.Errorf("unexpected end of NDJSON after index action")
			}

			var doc map[string]interface{}
			if err := json.Unmarshal([]byte(docLine), &doc); err != nil {
				return nil, err
			}
			op.Doc = doc

		case action.Delete != nil:
			op.Action = elasticsearch.BulkActionDelete
			op.Index = action.Delete.Index
			op.ID = action.Delete.ID

		default:
			continue
		}

		operations = append(operations, op)
	}

	return operations, nil
}
