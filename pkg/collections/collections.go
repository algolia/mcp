package collections

import "github.com/mark3labs/mcp-go/server"

// RegisterTools aggregates all collections tool registrations.
func RegisterTools(mcps *server.MCPServer) {
	RegisterListCollections(mcps)
	RegisterGetCollection(mcps)
	RegisterUpsertCollection(mcps)
	RegisterDeleteCollection(mcps)
	RegisterCommitCollection(mcps)
}
