package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterCopy(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	copyIndexTool := mcp.NewTool(
		"copy_index",
		mcp.WithDescription("Copy an index to a another index"),
		mcp.WithString(
			"indexName",
			mcp.Description("The name of the destination index"),
			mcp.Required(),
		),
	)

	mcps.AddTool(copyIndexTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dst, ok := req.Params.Arguments["indexName"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid indexName format, expected JSON string"), nil
		}
		res, err := client.OperationIndex(client.NewApiOperationIndexRequest(indexName, &search.OperationIndexParams{
			Operation:   "copy",
			Destination: dst,
		}))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not copy index: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("task", res)
	})
}
