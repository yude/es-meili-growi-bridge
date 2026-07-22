package translator_test

import (
	"testing"

	es "es-meili-growi-bridge/internal/elasticsearch"
	"es-meili-growi-bridge/internal/translator"
)

func TestBulkTranslator_IndexOperations(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{
		Operations: []es.BulkOperation{
			{
				Action: es.BulkActionIndex,
				Index:  "crowi",
				ID:     "doc1",
				Doc:    map[string]interface{}{"title": "Hello", "body": "World"},
			},
			{
				Action: es.BulkActionIndex,
				Index:  "crowi",
				ID:     "doc2",
				Doc:    map[string]interface{}{"title": "Foo", "body": "Bar"},
			},
		},
	}

	docs, deletes, indexName := trans.ToMeilisearchOperations(req)

	if len(docs) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(docs))
	}
	if len(deletes) != 0 {
		t.Errorf("expected 0 deletes, got %d", len(deletes))
	}
	if indexName != "crowi" {
		t.Errorf("expected indexName='crowi', got '%s'", indexName)
	}

	if docs[0]["id"] != "doc1" {
		t.Errorf("expected doc[0].id='doc1', got '%v'", docs[0]["id"])
	}
	if docs[0]["title"] != "Hello" {
		t.Errorf("expected doc[0].title='Hello', got '%v'", docs[0]["title"])
	}
}

func TestBulkTranslator_DeleteOperations(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{
		Operations: []es.BulkOperation{
			{
				Action: es.BulkActionDelete,
				Index:  "crowi",
				ID:     "doc1",
			},
			{
				Action: es.BulkActionDelete,
				Index:  "crowi",
				ID:     "doc2",
			},
		},
	}

	docs, delIDs, _ := trans.ToMeilisearchOperations(req)

	if len(docs) != 0 {
		t.Errorf("expected 0 documents, got %d", len(docs))
	}
	if len(delIDs) != 2 {
		t.Fatalf("expected 2 deletes, got %d", len(delIDs))
	}
	if delIDs[0] != "doc1" {
		t.Errorf("expected deletes[0]='doc1', got '%s'", delIDs[0])
	}
}

func TestBulkTranslator_MixedOperations(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{
		Operations: []es.BulkOperation{
			{
				Action: es.BulkActionIndex,
				Index:  "crowi",
				ID:     "doc1",
				Doc:    map[string]interface{}{"title": "Test"},
			},
			{
				Action: es.BulkActionDelete,
				Index:  "crowi",
				ID:     "old_doc",
			},
		},
	}

	docs, deletes, _ := trans.ToMeilisearchOperations(req)

	if len(docs) != 1 {
		t.Errorf("expected 1 document, got %d", len(docs))
	}
	if len(deletes) != 1 {
		t.Errorf("expected 1 delete, got %d", len(deletes))
	}
}

func TestBulkTranslator_EmptyOperations(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{Operations: []es.BulkOperation{}}

	docs, deletes, _ := trans.ToMeilisearchOperations(req)

	if len(docs) != 0 {
		t.Errorf("expected 0 documents, got %d", len(docs))
	}
	if len(deletes) != 0 {
		t.Errorf("expected 0 deletes, got %d", len(deletes))
	}
}

func TestBulkTranslator_UpdateOperation(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{
		Operations: []es.BulkOperation{
			{
				Action: es.BulkActionUpdate,
				Index:  "crowi",
				ID:     "doc1",
				Doc:    map[string]interface{}{"title": "Updated"},
			},
		},
	}

	docs, _, _ := trans.ToMeilisearchOperations(req)

	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
}

func TestBulkTranslator_MissingID(t *testing.T) {
	trans := translator.NewBulkTranslator()
	req := &es.BulkRequestBody{
		Operations: []es.BulkOperation{
			{
				Action: es.BulkActionIndex,
				Index:  "crowi",
				Doc:    map[string]interface{}{"title": "No ID doc"},
			},
		},
	}

	docs, _, _ := trans.ToMeilisearchOperations(req)

	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
	if _, ok := docs[0]["id"]; ok {
		t.Errorf("did not expect 'id' field for document without ID")
	}
}
