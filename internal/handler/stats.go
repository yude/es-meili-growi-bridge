package handler

import (
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

func (h *Handlers) IndexStats(w http.ResponseWriter, r *http.Request) {
	esIndex := r.PathValue("index")
	index := meiliIndex(h, esIndex)

	stats, err := h.meiliClient.GetIndexStats(index)
	if err != nil {
		writeError(w, elasticsearch.NewNotFound(index))
		return
	}

	esStats := elasticsearch.IndexStats{
		UUID: esIndex,
		Primaries: elasticsearch.IndexPrimariesStats{
			Docs:     elasticsearch.DocsStats{Count: stats.NumberOfDocuments, Deleted: 0},
			Store:    elasticsearch.StoreStats{SizeInBytes: 0},
			Indexing: elasticsearch.IndexingStats{},
		},
		Total: elasticsearch.IndexPrimariesStats{
			Docs:     elasticsearch.DocsStats{Count: stats.NumberOfDocuments, Deleted: 0},
			Store:    elasticsearch.StoreStats{SizeInBytes: 0},
			Indexing: elasticsearch.IndexingStats{},
		},
	}

	resp := &elasticsearch.IndicesStatsResponse{
		Shards:  elasticsearch.DefaultShards(),
		All: map[string]interface{}{
			"primaries": map[string]interface{}{
				"docs": elasticsearch.DocsStats{Count: stats.NumberOfDocuments, Deleted: 0},
			},
			"total": map[string]interface{}{
				"docs": elasticsearch.DocsStats{Count: stats.NumberOfDocuments, Deleted: 0},
			},
		},
		Indices: map[string]elasticsearch.IndexStats{
			esIndex: esStats,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}
