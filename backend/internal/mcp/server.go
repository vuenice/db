package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"chatdb/internal/engine"
)

// JSONRPCRequest represents an incoming JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents an outgoing JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id,omitempty"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError holds JSON-RPC error details.
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ToolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// List of available MCP tools.
func availableTools() []map[string]any {
	return []map[string]any{
		{
			"name":        "list_tables",
			"description": "List all tables and views in the database. Defaults to 'public' schema for PostgreSQL.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"schema": map[string]any{
						"type":        "string",
						"description": "Schema name (optional, defaults to public on PostgreSQL)",
					},
				},
			},
		},
		{
			"name":        "describe_table",
			"description": "Get columns, data types, nullability, and indexes for a specific table.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"table": map[string]any{
						"type":        "string",
						"description": "Table name",
					},
					"schema": map[string]any{
						"type":        "string",
						"description": "Schema name (optional, defaults to public on PostgreSQL)",
					},
				},
				"required": []string{"table"},
			},
		},
		{
			"name":        "read_query",
			"description": "Execute a safe read-only SQL query (SELECT, SHOW, EXPLAIN) and return tabular results.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"sql": map[string]any{
						"type":        "string",
						"description": "SQL statement (SELECT only)",
					},
					"max_rows": map[string]any{
						"type":        "integer",
						"description": "Maximum rows to return (default 100, max 1000)",
					},
				},
				"required": []string{"sql"},
			},
		},
		{
			"name":        "list_databases",
			"description": "List all physical databases available on the server.",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

// HandleMessage processes a single JSON-RPC 2.0 request and returns the serialized response.
func HandleMessage(ctx context.Context, eng engine.Engine, defaultDB string, rawReq []byte) ([]byte, error) {
	var req JSONRPCRequest
	if err := json.Unmarshal(rawReq, &req); err != nil {
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &JSONRPCError{Code: -32700, Message: "Parse error"},
		}
		return json.Marshal(resp)
	}

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "vuenicedb-mcp",
				"version": "1.0.0",
			},
		}

	case "notifications/initialized":
		return nil, nil // No response for notifications

	case "ping":
		resp.Result = map[string]any{}

	case "tools/list":
		resp.Result = map[string]any{
			"tools": availableTools(),
		}

	case "tools/call":
		var p ToolCallParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			resp.Error = &JSONRPCError{Code: -32602, Message: "Invalid params"}
			break
		}

		resultText, err := executeTool(ctx, eng, p.Name, p.Arguments)
		if err != nil {
			resp.Result = map[string]any{
				"isError": true,
				"content": []map[string]any{
					{
						"type": "text",
						"text": fmt.Sprintf("Error: %v", err),
					},
				},
			}
		} else {
			resp.Result = map[string]any{
				"content": []map[string]any{
					{
						"type": "text",
						"text": resultText,
					},
				},
			}
		}

	default:
		resp.Error = &JSONRPCError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)}
	}

	return json.Marshal(resp)
}

func executeTool(ctx context.Context, eng engine.Engine, name string, args map[string]any) (string, error) {
	if eng == nil {
		return "", errors.New("no database connection available")
	}

	switch name {
	case "list_tables":
		schema := "public"
		if s, ok := args["schema"].(string); ok && strings.TrimSpace(s) != "" {
			schema = strings.TrimSpace(s)
		}
		tables, err := eng.ListTables(ctx, schema)
		if err != nil {
			return "", err
		}
		data, _ := json.MarshalIndent(tables, "", "  ")
		return string(data), nil

	case "describe_table":
		table, ok := args["table"].(string)
		if !ok || strings.TrimSpace(table) == "" {
			return "", errors.New("table parameter is required")
		}
		table = strings.TrimSpace(table)
		schema := "public"
		if s, ok := args["schema"].(string); ok && strings.TrimSpace(s) != "" {
			schema = strings.TrimSpace(s)
		}
		cols, err := eng.ListColumns(ctx, schema, table)
		if err != nil {
			return "", fmt.Errorf("list columns: %w", err)
		}
		idx, _ := eng.ListIndexes(ctx, schema, table)
		res := map[string]any{
			"schema":  schema,
			"table":   table,
			"columns": cols,
			"indexes": idx,
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return string(data), nil

	case "read_query":
		sqlStr, ok := args["sql"].(string)
		if !ok || strings.TrimSpace(sqlStr) == "" {
			return "", errors.New("sql parameter is required")
		}
		sqlStr = strings.TrimSpace(sqlStr)

		// Ensure read-only query
		upper := strings.ToUpper(sqlStr)
		if strings.HasPrefix(upper, "INSERT") || strings.HasPrefix(upper, "UPDATE") ||
			strings.HasPrefix(upper, "DELETE") || strings.HasPrefix(upper, "DROP") ||
			strings.HasPrefix(upper, "ALTER") || strings.HasPrefix(upper, "TRUNCATE") {
			return "", errors.New("only read queries (SELECT, SHOW, EXPLAIN) are permitted via MCP")
		}

		maxRows := 100
		if mr, ok := args["max_rows"].(float64); ok && mr > 0 {
			maxRows = int(mr)
			if maxRows > 1000 {
				maxRows = 1000
			}
		}

		qr, err := eng.Execute(ctx, sqlStr, maxRows)
		if err != nil {
			return "", err
		}
		data, _ := json.MarshalIndent(qr, "", "  ")
		return string(data), nil

	case "list_databases":
		dbs, err := eng.ListDatabases(ctx)
		if err != nil {
			return "", err
		}
		data, _ := json.MarshalIndent(dbs, "", "  ")
		return string(data), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// RunStdio starts a JSON-RPC 2.0 loop reading from stdin and writing to stdout.
func RunStdio(ctx context.Context, eng engine.Engine, defaultDB string) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		trimmed := strings.TrimSpace(string(line))
		if trimmed == "" {
			continue
		}

		reqCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		respBytes, err := HandleMessage(reqCtx, eng, defaultDB, []byte(trimmed))
		cancel()

		if err != nil {
			continue
		}
		if len(respBytes) > 0 {
			os.Stdout.Write(respBytes)
			os.Stdout.WriteString("\n")
		}
	}
}
