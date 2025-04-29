package records

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterGetObject(mcps *server.MCPServer, index *search.Index) {
	getObjectTool := mcp.NewTool(
		"get_object",
		mcp.WithDescription("Get an object by its object ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The object ID to look up"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, _ := req.Params.Arguments["objectID"].(string)

		var x map[string]any
		if err := index.GetObject(objectID, &x); err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not get object: %v", err),
			), nil
		}
		return mcputil.JSONToolResult("object", x)
	})
}
