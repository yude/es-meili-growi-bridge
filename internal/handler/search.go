package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")
	actualIndex := meiliIndex(h, resolveIndex(h, index))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to read request body"))
		return
	}

	if len(body) == 0 {
		writeError(w, elasticsearch.NewBadRequest("request body is required"))
		return
	}

	var req elasticsearch.SearchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("[search] invalid JSON for index=%s body=%s err=%s", actualIndex, string(body), err.Error())
		writeError(w, elasticsearch.NewBadRequest("invalid JSON: "+err.Error()))
		return
	}

	log.Printf("[search] index=%s body=%s", actualIndex, string(body))

	meiliParams, err := h.searchTrans.ToMeilisearchParams(&req, actualIndex)
	if err != nil {
		log.Printf("[search] translation error for index=%s err=%s", actualIndex, err.Error())
		writeError(w, elasticsearch.NewBadRequest("translation error: "+err.Error()))
		return
	}

	if meiliParams.Q == "" {
		meiliParams.Q = ""
	}

	result, err := h.meiliClient.Search(actualIndex, meiliParams)
	if err != nil {
		log.Printf("[search] search error for index=%s err=%s", actualIndex, err.Error())
		writeError(w, elasticsearch.NewBadRequest("search error: "+err.Error()))
		return
	}

	esResp := h.responseTrans.ToESSearchResponse(result, actualIndex)
	writeJSON(w, http.StatusOK, esResp)
}
