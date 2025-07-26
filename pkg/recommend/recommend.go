package recommend

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all Recommend tools with the MCP server.
func RegisterAll(mcps *mcp.Server) {
	// Register all Recommend tools.
	RegisterGetRecommendations(mcps)
	RegisterGetRecommendRule(mcps)
	RegisterDeleteRecommendRule(mcps)
	RegisterSearchRecommendRules(mcps)
	RegisterBatchRecommendRules(mcps)
	RegisterGetRecommendStatus(mcps)
}
