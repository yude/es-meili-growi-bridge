package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"es-meili-growi-bridge/internal/config"
	"es-meili-growi-bridge/internal/elasticsearch"
	"es-meili-growi-bridge/internal/meilisearch"
	"es-meili-growi-bridge/internal/translator"
)

type Handlers struct {
	meiliClient *meilisearch.Client
	cfg         *config.Config
	aliasStore  *AliasStore

	searchTrans   translator.SearchTranslator
	bulkTrans     translator.BulkTranslator
	responseTrans translator.ResponseTranslator

	indexMappingStore map[string]map[string]interface{}
}

func NewHandlers(cfg *config.Config, meiliClient *meilisearch.Client) *Handlers {
	return &Handlers{
		meiliClient: meiliClient,
		cfg:         cfg,
		aliasStore:  NewAliasStore(),

		searchTrans:   translator.NewSearchTranslator(),
		bulkTrans:     translator.NewBulkTranslator(),
		responseTrans: translator.NewResponseTranslator(),

		indexMappingStore: make(map[string]map[string]interface{}),
	}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Root)
	mux.HandleFunc("GET /_cluster/health", h.ClusterHealth)
	mux.HandleFunc("GET /_nodes/info", h.NodesInfo)
	mux.HandleFunc("GET /_cat/indices", h.CatIndices)
	mux.HandleFunc("GET /_cat/aliases", h.CatAliases)
	mux.HandleFunc("POST /_bulk", h.Bulk)
	mux.HandleFunc("POST /_reindex", h.Reindex)
	mux.HandleFunc("POST /_aliases", h.UpdateAliases)
	mux.HandleFunc("HEAD /{index}", h.IndexExists)
	mux.HandleFunc("PUT /{index}", h.CreateIndex)
	mux.HandleFunc("DELETE /{index}", h.DeleteIndex)
	mux.HandleFunc("GET /{index}/_alias", h.GetAliases)
	mux.HandleFunc("HEAD /{index}/_alias/{name}", h.AliasExists)
	mux.HandleFunc("PUT /{index}/_alias/{name}", h.PutAlias)
	mux.HandleFunc("GET /{index}/_stats", h.IndexStats)
	mux.HandleFunc("POST /{index}/_search", h.Search)
	mux.HandleFunc("GET /{index}/_validate/query", h.ValidateQuery)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, err *elasticsearch.Error) {
	body := elasticsearch.ErrorResponse{
		Error: elasticsearch.ErrorBody{
			Type:   err.ErrorType,
			Reason: err.Reason,
			Index:  err.Index,
		},
		Status: err.Status,
	}
	writeJSON(w, err.Status, body)
}

func writeNotImplemented(w http.ResponseWriter) {
	writeError(w, elasticsearch.NewNotImplemented())
}

func resolveIndex(h *Handlers, index string) string {
	if resolved, ok := h.aliasStore.ResolveAlias(index); ok {
		return resolved
	}
	return index
}
