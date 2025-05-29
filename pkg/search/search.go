package search

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/search/indices"
	"github.com/algolia/mcp/pkg/search/query"
	"github.com/algolia/mcp/pkg/search/records"
	"github.com/mark3labs/mcp-go/server"
	"github.com/algolia/mcp/pkg/mcputil/client"
)

// RegisterAll registers all Search tools with the MCP server (both read and write).
func RegisterAll(mcps *server.MCPServer) {
	var searchClient *search.Client
	var searchIndex *search.Index

	searchClient, searchIndex := client.clientInstance()

	// Register both read and write operations.
	RegisterReadAll(mcps, searchClient, searchIndex)
	RegisterWriteAll(mcps, searchClient, searchIndex)
}

// RegisterReadAll registers read-only Search tools with the MCP server.
func RegisterReadAll(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	// Register read-only operations.
	indices.RegisterList(mcps, client)
	indices.RegisterGetSettings(mcps, index)
	query.RegisterRunQuery(mcps, client, index)
	records.RegisterGetObject(mcps, index)
}

// RegisterWriteAll registers write-only Search tools with the MCP server.
func RegisterWriteAll(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	// Register write operations.
	indices.RegisterClear(mcps, index)
	indices.RegisterCopy(mcps, client, index)
	indices.RegisterDelete(mcps, index)
	indices.RegisterMove(mcps, client, index)
	indices.RegisterSetSettings(mcps, index)
	records.RegisterDeleteObject(mcps, index)
	records.RegisterInsertObject(mcps, index)
	records.RegisterInsertObjects(mcps, index)
}
