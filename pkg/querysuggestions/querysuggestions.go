package querysuggestions

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all Query Suggestions tools with the MCP server.
func RegisterAll(mcps *mcp.Server) {
	// Register all Query Suggestions tools.
	RegisterListConfigs(mcps)
	RegisterGetConfig(mcps)
	RegisterCreateConfig(mcps)
	RegisterUpdateConfig(mcps)
	RegisterDeleteConfig(mcps)
	RegisterGetConfigStatus(mcps)
	RegisterGetLogFile(mcps)
}
