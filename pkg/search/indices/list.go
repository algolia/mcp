package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterList(mcps *server.MCPServer, client *search.Client) {
	listIndexTool := mcp.NewTool(
		"list_indices",
		mcp.WithDescription("List the indices in the application"),
	)

	mcps.AddTool(listIndexTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := client.ListIndices()
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not list indices: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("indices", res)
	})
}
