package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterClear(mcps *server.MCPServer, index *search.Index) {
	clearIndexTool := mcp.NewTool(
		"clear_index",
		mcp.WithDescription("Clear an index by removing all records"),
	)

	mcps.AddTool(clearIndexTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := index.ClearObjects()
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not clear index: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("object", res)
	})
}
