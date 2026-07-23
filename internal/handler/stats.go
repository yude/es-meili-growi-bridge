package handler

import (
	"net/http"
	"strings"

	"es-meili-growi-bridge/internal/elasticsearch"
	"es-meili-growi-bridge/internal/meilisearch"
)

func (h *Handlers) IndexStats(w http.ResponseWriter, r *http.Request) {
	esIndices := strings.Split(r.PathValue("index"), ",")

	indices := make(map[string]elasticsearch.IndexStats, len(esIndices))
	var totalDocs int
	var meiliStats *meilisearch.IndexStats

	for _, esIdx := range esIndices {
		meiliIdx := meiliIndex(h, esIdx)
		if meiliStats == nil {
			s, err := h.meiliClient.GetIndexStats(meiliIdx)
			if err != nil {
				writeError(w, elasticsearch.NewNotFound(meiliIdx))
				return
			}
			meiliStats = s
		}
		esStats := elasticsearch.IndexStats{
			UUID: esIdx,
			Primaries: elasticsearch.IndexPrimariesStats{
				Docs:  elasticsearch.DocsStats{Count: meiliStats.NumberOfDocuments, Deleted: 0},
				Store: elasticsearch.StoreStats{SizeInBytes: 0},
			},
			Total: elasticsearch.IndexPrimariesStats{
				Docs:  elasticsearch.DocsStats{Count: meiliStats.NumberOfDocuments, Deleted: 0},
				Store: elasticsearch.StoreStats{SizeInBytes: 0},
			},
		}
		totalDocs += int(meiliStats.NumberOfDocuments)
		indices[esIdx] = esStats
	}

	resp := &elasticsearch.IndicesStatsResponse{
		Shards: elasticsearch.DefaultShards(),
		All: map[string]interface{}{
			"primaries": elasticsearch.IndexPrimariesStats{
				Docs: elasticsearch.DocsStats{Count: totalDocs, Deleted: 0},
			},
			"total": elasticsearch.IndexPrimariesStats{
				Docs: elasticsearch.DocsStats{Count: totalDocs, Deleted: 0},
			},
		},
		Indices: indices,
	}

	writeJSON(w, http.StatusOK, resp)
}
