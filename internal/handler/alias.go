package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) GetAliases(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")

	aliases := h.aliasStore.GetIndexAliases(index)

	resp := map[string]*elasticsearch.IndexAliasInfo{
		index: {Aliases: aliases},
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) AliasExists(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")
	alias := r.PathValue("name")

	if h.aliasStore.ExistsAlias(alias, index) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *Handlers) PutAlias(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")
	alias := r.PathValue("name")

	h.aliasStore.PutAlias(index, alias)

	writeJSON(w, http.StatusOK, elasticsearch.NewPutAliasResponse())
}

func (h *Handlers) UpdateAliases(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to read request body"))
		return
	}

	var req elasticsearch.UpdateAliasesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, elasticsearch.NewBadRequest("invalid JSON: "+err.Error()))
		return
	}

	for _, action := range req.Actions {
		if action.Add != nil {
			h.aliasStore.PutAlias(action.Add.Index, action.Add.Alias)
		}
		if action.Remove != nil {
			h.aliasStore.RemoveAlias(action.Remove.Index, action.Remove.Alias)
		}
	}

	writeJSON(w, http.StatusOK, elasticsearch.NewUpdateAliasesResponse())
}
