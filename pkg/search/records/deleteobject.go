package records

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterDeleteObject(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	deleteObjectTool := mcp.NewTool(
		"delete_object",
		mcp.WithDescription("Delete an object by its object ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The object ID to delete"),
			mcp.Required(),
		),
	)

	mcps.AddTool(deleteObjectTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, _ := req.Params.Arguments["objectID"].(string)

		res, err := client.DeleteObject(client.NewApiDeleteObjectRequest(indexName, objectID))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not delete object: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("object", res)
	})
}
