package recommend

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Recommend tools with the MCP server.
func RegisterAll(mcps *server.MCPServer) {
	// Register all Recommend tools.
	RegisterGetRecommendations(mcps)
	RegisterGetRecommendRule(mcps)
	RegisterDeleteRecommendRule(mcps)
	RegisterSearchRecommendRules(mcps)
	RegisterBatchRecommendRules(mcps)
	RegisterGetRecommendStatus(mcps)
}
