package monitoring

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterTools aggregates all monitoring tool registrations.
func RegisterTools(s *mcp.Server) {
	RegisterGetClusterIncidents(s)
	RegisterGetClusterStatus(s)
	RegisterGetClustersStatus(s)
	RegisterGetIncidents(s)
	RegisterGetIndexingTime(s)
	RegisterGetLatency(s)
	RegisterGetMetrics(s)
	RegisterGetReachability(s)
	RegisterGetServers(s)
}
