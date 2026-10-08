package api

import (
	"errors"
	"io"
	"net/http"

	"chatdb/internal/mcp"
)

// handleConnectionMCP handles Model Context Protocol JSON-RPC requests for a specific connection.
func (s *Server) handleConnectionMCP(w http.ResponseWriter, r *http.Request) {
	eng, conn, err := s.resolveEngine(r, false)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errors.New("read body failed"))
		return
	}

	resp, err := mcp.HandleMessage(r.Context(), eng, conn.Database, body)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
