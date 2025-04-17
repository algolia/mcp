package rules

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetRule(mcps *server.MCPServer, index *search.Index) {
	getRuleTool := mcp.NewTool(
		"get_rule",
		mcp.WithDescription("Get a rule from the Algolia index by its ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The unique identifier of the rule to retrieve"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getRuleTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objectID format"), nil
		}

		rule, err := index.GetRule(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not get rule: %w", err)
		}

		return mcputil.JSONToolResult("rule", rule)
	})
}
