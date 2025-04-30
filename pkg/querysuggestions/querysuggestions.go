package querysuggestions

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Query Suggestions tools with the MCP server.
func RegisterAll(mcps *server.MCPServer) {
	// Register all Query Suggestions tools.
	RegisterListConfigs(mcps)
	RegisterGetConfig(mcps)
	RegisterCreateConfig(mcps)
	RegisterUpdateConfig(mcps)
	RegisterDeleteConfig(mcps)
	RegisterGetConfigStatus(mcps)
	RegisterGetLogFile(mcps)
}
