package search

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/search/indices"
	"github.com/algolia/mcp/pkg/search/query"
	"github.com/algolia/mcp/pkg/search/records"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Search tools with the MCP server (both read and write).
func RegisterAll(mcps *server.MCPServer) {
	// Initialize Algolia client.
	// Note: In a real implementation, you would get the app ID and API key from environment variables.
	client := search.NewClient("", "")
	index := client.InitIndex("default_index")

	// Register both read and write operations.
	RegisterReadAll(mcps, client, index)
	RegisterWriteAll(mcps, client, index)
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
