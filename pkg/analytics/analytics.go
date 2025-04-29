package analytics

import "github.com/mark3labs/mcp-go/server"

// RegisterTools aggregates all analytics tool registrations.
func RegisterTools(mcps *server.MCPServer) {
	RegisterGetClickThroughRate(mcps)
	RegisterGetNoResultsRate(mcps)
	RegisterGetSearchesCount(mcps)
	RegisterGetTopSearches(mcps)
}
