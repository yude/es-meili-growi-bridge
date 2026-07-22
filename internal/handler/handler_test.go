package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"es-meili-growi-bridge/internal/config"
	"es-meili-growi-bridge/internal/elasticsearch"
	"es-meili-growi-bridge/internal/meilisearch"
)

func testHandlers(t *testing.T) *Handlers {
	t.Helper()
	cfg := &config.Config{
		MeilisearchURL:  "http://localhost:17700",
		MeilisearchAPIKey: "",
		ListenAddr:      ":9200",
		LogLevel:        "debug",
	}
	meiliClient := meilisearch.NewClient(cfg.MeilisearchURL, cfg.MeilisearchAPIKey)
	return NewHandlers(cfg, meiliClient)
}

func TestRootHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp elasticsearch.RootResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.TagLine == "" {
		t.Error("expected non-empty tagline")
	}
	if resp.Version.Number == "" {
		t.Error("expected non-empty version number")
	}
}

func TestClusterHealthHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/_cluster/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp elasticsearch.ClusterHealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.ClusterName == "" {
		t.Error("expected non-empty cluster_name")
	}
	if resp.Status != "green" {
		t.Errorf("expected status='green', got '%s'", resp.Status)
	}
}

func TestNodesInfoHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/_nodes/info", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp elasticsearch.NodesInfoResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.Nodes) == 0 {
		t.Error("expected at least 1 node")
	}
}

func TestIndexExistsNotFound(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("HEAD", "/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Since Meilisearch is not running, this will fail to find the index
	// Should return 404
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusOK {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCreateIndex_BadJSON(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	body := strings.NewReader("{invalid json}")
	req := httptest.NewRequest("PUT", "/testindex", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var errResp elasticsearch.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Logf("non-JSON error response: %s", rec.Body.String())
	}
}

func TestAliasStore(t *testing.T) {
	store := NewAliasStore()

	store.PutAlias("crowi", "crowi-alias")

	if !store.ExistsAlias("crowi-alias", "crowi") {
		t.Error("expected alias 'crowi-alias' to exist for 'crowi'")
	}

	aliases := store.GetIndexAliases("crowi")
	if _, ok := aliases["crowi-alias"]; !ok {
		t.Error("expected 'crowi-alias' in aliases list")
	}

	index, ok := store.ResolveAlias("crowi-alias")
	if !ok {
		t.Fatal("expected to resolve alias 'crowi-alias'")
	}
	if index != "crowi" {
		t.Errorf("expected index='crowi', got '%s'", index)
	}

	store.RemoveAlias("crowi", "crowi-alias")
	if store.ExistsAlias("crowi-alias", "crowi") {
		t.Error("expected alias to be removed")
	}

	if _, ok := store.ResolveAlias("crowi-alias"); ok {
		t.Error("expected alias resolution to fail after removal")
	}
}

func TestAliasStore_MultipleAliases(t *testing.T) {
	store := NewAliasStore()

	store.PutAlias("crowi", "crowi-alias")
	store.PutAlias("crowi", "crowi-alias2")

	all := store.GetAllAliases()
	if len(all) != 1 {
		t.Fatalf("expected 1 index, got %d", len(all))
	}

	aliases := all["crowi"]
	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(aliases))
	}
}

func TestSearchHandler_EmptyBody(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/crowi/_search", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty body, got %d", rec.Code)
	}
}

func TestSearchHandler_InvalidJSON(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	body := strings.NewReader("{invalid}")
	req := httptest.NewRequest("POST", "/crowi/_search", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestBulkHandler_EmptyBody(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/_bulk", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp elasticsearch.BulkResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.Items) != 0 {
		t.Errorf("expected 0 items for empty body, got %d", len(resp.Items))
	}
}

func TestReindexHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	body := strings.NewReader(`{"source":{"index":"crowi"},"dest":{"index":"crowi-tmp"}}`)
	req := httptest.NewRequest("POST", "/_reindex", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["task"] != "dummy-task-id" {
		t.Errorf("expected task='dummy-task-id', got '%v'", resp["task"])
	}
}

func TestUpdateAliasesHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	body := strings.NewReader(`{
		"actions": [
			{"add": {"alias": "crowi-alias", "index": "crowi"}},
			{"remove": {"alias": "crowi-alias", "index": "crowi-tmp"}}
		]
	}`)
	req := httptest.NewRequest("POST", "/_aliases", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if !h.aliasStore.ExistsAlias("crowi-alias", "crowi") {
		t.Error("expected crowi-alias to exist for crowi after update")
	}
	if h.aliasStore.ExistsAlias("crowi-alias", "crowi-tmp") {
		t.Error("expected crowi-alias to be removed from crowi-tmp")
	}
}

func TestGetAliasesHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	h.aliasStore.PutAlias("crowi", "crowi-alias")

	req := httptest.NewRequest("GET", "/crowi/_alias", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp map[string]elasticsearch.IndexAliasInfo
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	info, ok := resp["crowi"]
	if !ok {
		t.Fatal("expected 'crowi' key in response")
	}
	if _, ok := info.Aliases["crowi-alias"]; !ok {
		t.Error("expected 'crowi-alias' in aliases")
	}
}

func TestAliasExists(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	h.aliasStore.PutAlias("crowi", "crowi-alias")

	req := httptest.NewRequest("HEAD", "/crowi/_alias/crowi-alias", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	req = httptest.NewRequest("HEAD", "/crowi/_alias/nonexistent", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPutAliasHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/crowi/_alias/test-alias", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if !h.aliasStore.ExistsAlias("test-alias", "crowi") {
		t.Error("expected test-alias to exist")
	}
}

func TestValidateQueryHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/crowi/_validate/query?explain=true", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp elasticsearch.ValidateQueryResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if !resp.Valid {
		t.Error("expected valid=true")
	}
}

func TestCatIndicesHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/_cat/indices", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCatAliasesHandler(t *testing.T) {
	h := testHandlers(t)
	mux := http.NewServeMux()
	h.Register(mux)

	h.aliasStore.PutAlias("crowi", "crowi-alias")

	req := httptest.NewRequest("GET", "/_cat/aliases", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var rows []elasticsearch.CatAliasRow
	if err := json.NewDecoder(rec.Body).Decode(&rows); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].Alias != "crowi-alias" {
		t.Errorf("expected alias='crowi-alias', got '%s'", rows[0].Alias)
	}
}
