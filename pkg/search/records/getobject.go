package records

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterGetObject(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	getObjectTool := mcp.NewTool(
		"get_object",
		mcp.WithDescription("Get an object by its object ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The object ID to look up"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getObjectTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, _ := req.Params.Arguments["objectID"].(string)

		x, err := client.GetObject(client.NewApiGetObjectRequest(indexName, objectID))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not get object: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("object", x)
	})
}
