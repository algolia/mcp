package analytics

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterTools aggregates all analytics tool registrations.
func RegisterTools(mcps *mcp.Server) {
	RegisterGetClickThroughRate(mcps)
	RegisterGetNoResultsRate(mcps)
	RegisterGetSearchesCount(mcps)
	RegisterGetTopSearches(mcps)
}
