package rules

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterClearRules(mcps *server.MCPServer, writeIndex *search.Index) {
	clearRulesTool := mcp.NewTool(
		"clear_rules",
		mcp.WithDescription("Clear all rules from the Algolia index"),
	)

	mcps.AddTool(clearRulesTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if writeIndex == nil {
			return mcp.NewToolResultError("write API key not set, cannot clear rules"), nil
		}

		res, err := writeIndex.ClearRules()
		if err != nil {
			return nil, fmt.Errorf("could not clear rules: %w", err)
		}

		return mcputil.JSONToolResult("clear result", res)
	})
}
