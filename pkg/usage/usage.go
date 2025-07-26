package usage

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all Usage tools with the MCP server.
func RegisterAll(mcps *mcp.Server) {
	// Register all Usage tools.
	RegisterGetMetricsRegistry(mcps)
	RegisterGetDailyMetrics(mcps)
	RegisterGetHourlyMetrics(mcps)
}
