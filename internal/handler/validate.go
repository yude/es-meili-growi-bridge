package handler

import (
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) ValidateQuery(w http.ResponseWriter, r *http.Request) {
	index := r.PathValue("index")

	resp := elasticsearch.NewValidateQueryExplanation(index, true)
	writeJSON(w, http.StatusOK, resp)
}
