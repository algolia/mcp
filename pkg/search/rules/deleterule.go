package rules

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterDeleteRule(mcps *server.MCPServer, index *search.Index) {
	deleteRuleTool := mcp.NewTool(
		"delete_rule",
		mcp.WithDescription("Delete a rule by its object ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The object ID to delete"),
			mcp.Required(),
		),
	)

	mcps.AddTool(deleteRuleTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid object format, expected JSON string"), nil
		}

		resp, err := index.DeleteRule(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not delete rule: %w", err)
		}

		return mcputil.JSONToolResult("rule", resp)
	})
}
