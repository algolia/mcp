package abtesting

import "github.com/mark3labs/mcp-go/server"

// RegisterTools aggregates all abtesting tool registrations.
func RegisterTools(mcps *server.MCPServer) {
	RegisterListABTests(mcps)
	RegisterGetABTest(mcps)
	RegisterCreateABTest(mcps)
	RegisterDeleteABTest(mcps)
	RegisterStopABTest(mcps)
	RegisterEstimateABTest(mcps)
	RegisterScheduleABTest(mcps)
}
