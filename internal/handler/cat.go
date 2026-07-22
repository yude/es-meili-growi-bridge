package handler

import (
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) CatIndices(w http.ResponseWriter, r *http.Request) {
	indexes, err := h.meiliClient.ListIndexes()
	if err != nil {
		writeJSON(w, http.StatusOK, []elasticsearch.CatIndexRow{})
		return
	}

	rows := make([]elasticsearch.CatIndexRow, 0, len(indexes))
	for _, idx := range indexes {
		rows = append(rows, elasticsearch.CatIndexRow{
			Health:   "green",
			Status:   "open",
			Index:    idx.UID,
			UUID:     idx.UID,
			Pri:      "1",
			Rep:      "0",
			DocsCount: "0",
		})
	}

	writeJSON(w, http.StatusOK, rows)
}

func (h *Handlers) CatAliases(w http.ResponseWriter, r *http.Request) {
	allAliases := h.aliasStore.GetAllAliases()

	rows := make([]elasticsearch.CatAliasRow, 0)
	for index, aliases := range allAliases {
		for alias := range aliases {
			rows = append(rows, elasticsearch.CatAliasRow{
				Alias:  alias,
				Index:  index,
				Filter: "-",
			})
		}
	}

	writeJSON(w, http.StatusOK, rows)
}
