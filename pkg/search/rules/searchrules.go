package rules

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterSearchRules(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	searchRulesTool := mcp.NewTool(
		"search_rules",
		mcp.WithDescription("Search for rules in the Algolia index"),
		mcp.WithString(
			"query",
			mcp.Description("The query to search for"),
			mcp.Required(),
		),
		mcp.WithString(
			"anchoring",
			mcp.Description("When specified, restricts matches to rules with a specific anchoring type. When omitted, all anchoring types may match."),
			mcp.Enum("is", "contains", "startsWith", "endsWith"),
		),
		mcp.WithString(
			"context",
			mcp.Description("When specified, restricts matches to contextual rules with a specific context. When omitted, all contexts may match."),
		),
		mcp.WithBoolean(
			"enabled",
			mcp.Description("When specified, restricts matches to rules with a specific enabled status. When omitted, all enabled statuses may match."),
		),
	)

	mcps.AddTool(searchRulesTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params search.SearchRulesParams
		if anchoring, ok := req.Params.Arguments["anchoring"].(string); ok {
			anchoringEnum, err := search.NewAnchoringFromValue(anchoring)
			if err != nil {
				return nil, fmt.Errorf("invalid anchoring: %w", err)
			}
			params.Anchoring = anchoringEnum
		}
		if context, ok := req.Params.Arguments["context"].(string); ok {
			params.Context = &context
		}
		if enabled, ok := req.Params.Arguments["enabled"].(bool); ok {
			params.Enabled.Set(&enabled)
		}
		query, _ := req.Params.Arguments["query"].(string)
		params.Query = &query

		resp, err := client.SearchRules(client.NewApiSearchRulesRequest(indexName).WithSearchRulesParams(&params))
		if err != nil {
			return nil, fmt.Errorf("could not search rules: %w", err)
		}

		return mcputil.JSONToolResult("rules", resp)
	})
}
