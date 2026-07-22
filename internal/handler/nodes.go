package handler

import (
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) NodesInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, elasticsearch.DefaultNodesInfo())
}
