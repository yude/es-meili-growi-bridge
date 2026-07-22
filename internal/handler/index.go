package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) IndexExists(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")

	_, err := h.meiliClient.GetIndex(index)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) CreateIndex(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to read request body"))
		return
	}

	var req elasticsearch.IndexCreateRequest
	if len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(w, elasticsearch.NewBadRequest("invalid JSON: "+err.Error()))
			return
		}
	}

	primaryKey := "id"
	if req.Mappings != nil {
		if props, ok := req.Mappings["properties"].(map[string]interface{}); ok {
			if idField, ok := props["id"].(map[string]interface{}); ok {
				if _, ok := idField["type"].(string); ok {
					primaryKey = "id"
				}
			}
		}
		h.indexMappingStore[index] = req.Mappings
	}

	_, err = h.meiliClient.CreateIndex(index, primaryKey)
	if err != nil {
		if esErr, ok := err.(*elasticsearch.Error); ok {
			writeError(w, esErr)
			return
		}
		writeError(w, elasticsearch.NewBadRequest("failed to create index: "+err.Error()))
		return
	}

	h.aliasStore.PutAlias(index, index+"-alias")

	resp := elasticsearch.NewIndexCreateResponse(index)
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) DeleteIndex(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")

	if err := h.meiliClient.DeleteIndex(index); err != nil {
		writeError(w, elasticsearch.NewNotFound(index))
		return
	}

	resp := elasticsearch.NewIndexDeleteResponse()
	writeJSON(w, http.StatusOK, resp)
}
