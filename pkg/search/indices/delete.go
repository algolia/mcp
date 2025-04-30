package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterDelete(mcps *server.MCPServer, index *search.Index) {
	deleteIndexTool := mcp.NewTool(
		"delete_index",
		mcp.WithDescription("Delete an index by removing all assets and configurations"),
	)

	mcps.AddTool(deleteIndexTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := index.Delete()
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not delete index: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("task", res)
	})
}
