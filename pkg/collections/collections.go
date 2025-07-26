package collections

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterTools aggregates all collections tool registrations.
func RegisterTools(mcps *mcp.Server) {
	RegisterListCollections(mcps)
	RegisterGetCollection(mcps)
	RegisterUpsertCollection(mcps)
	RegisterDeleteCollection(mcps)
	RegisterCommitCollection(mcps)
}
