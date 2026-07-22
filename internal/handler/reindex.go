package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"es-meili-growi-bridge/internal/elasticsearch"
)

type ReindexRequest struct {
	Source struct {
		Index string `json:"index"`
	} `json:"source"`
	Dest struct {
		Index string `json:"index"`
	} `json:"dest"`
	WaitForCompletion bool `json:"wait_for_completion,omitempty"`
}

type ReindexResponse struct {
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

func (h *Handlers) Reindex(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, elasticsearch.NewBadRequest("failed to read request body"))
		return
	}

	if len(body) == 0 {
		writeError(w, elasticsearch.NewBadRequest("request body is required"))
		return
	}

	var req ReindexRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, elasticsearch.NewBadRequest("invalid JSON: "+err.Error()))
		return
	}

	resp := ReindexResponse{
		Task:      "dummy-task-id",
		Completed: false,
	}

	if !req.WaitForCompletion {
		writeJSON(w, http.StatusOK, resp)
		return
	}

	resp.Completed = true
	writeJSON(w, http.StatusOK, resp)
}
