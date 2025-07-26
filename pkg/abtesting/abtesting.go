package abtesting

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterTools aggregates all abtesting tool registrations.
func RegisterTools(mcps *mcp.Server) {
	RegisterListABTests(mcps)
	RegisterGetABTest(mcps)
	RegisterCreateABTest(mcps)
	RegisterDeleteABTest(mcps)
	RegisterStopABTest(mcps)
	RegisterEstimateABTest(mcps)
	RegisterScheduleABTest(mcps)
}
