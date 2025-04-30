package usage

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Usage tools with the MCP server.
func RegisterAll(mcps *server.MCPServer) {
	// Register all Usage tools.
	RegisterGetMetricsRegistry(mcps)
	RegisterGetDailyMetrics(mcps)
	RegisterGetHourlyMetrics(mcps)
}
