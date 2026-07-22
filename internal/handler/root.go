package handler

import (
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) Root(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, elasticsearch.DefaultRootResponse())
}
