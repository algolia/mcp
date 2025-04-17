package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterClear(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	clearIndexTool := mcp.NewTool(
		"clear_index",
		mcp.WithDescription("Clear an index by removing all records"),
	)

	mcps.AddTool(clearIndexTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := client.ClearObjects(client.NewApiClearObjectsRequest(indexName))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not clear index: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("object", res)
	})
}
