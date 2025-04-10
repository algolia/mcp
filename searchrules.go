package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func RegisterSearchRules(mcps *server.MCPServer, index *search.Index) {
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

	mcps.AddTool(searchRulesTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.Params.Arguments["query"].(string)

		opts := []any{}
		if anchoring, ok := req.Params.Arguments["anchoring"].(string); ok {
			opts = append(opts, opt.Anchoring(anchoring))
		}
		if context, ok := req.Params.Arguments["context"].(string); ok {
			opts = append(opts, opt.RuleContexts(context))
		}
		if enabled, ok := req.Params.Arguments["enabled"].(bool); ok {
			opts = append(opts, opt.EnableRules(enabled))
		}

		resp, err := index.SearchRules(query, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not search rules: %w", err)
		}

		return jsonResponse("rules", resp)
	})
}
