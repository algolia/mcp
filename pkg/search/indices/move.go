package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterMove(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	moveIndexTool := mcp.NewTool(
		"move_index",
		mcp.WithDescription("Move an index to another index"),
		mcp.WithString(
			"indexName",
			mcp.Description("The name of the destination index"),
			mcp.Required(),
		),
	)

	mcps.AddTool(moveIndexTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dst, ok := req.Params.Arguments["indexName"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid indexName format, expected JSON string"), nil
		}

		res, err := client.CopyIndex(index.GetName(), dst)
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not move index: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("task", res)
	})
}
