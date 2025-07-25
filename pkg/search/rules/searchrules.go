package rules

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearchRulesParams defines the parameters for searching rules.
type SearchRulesParams struct {
	Query     string `json:"query" jsonschema:"The query to search for"`
	Anchoring string `json:"anchoring,omitempty" jsonschema:"When specified restricts matches to rules with a specific anchoring type. When omitted all anchoring types may match."`
	Context   string `json:"context,omitempty" jsonschema:"When specified restricts matches to contextual rules with a specific context. When omitted all contexts may match."`
	Enabled   *bool  `json:"enabled,omitempty" jsonschema:"When specified restricts matches to rules with a specific enabled status. When omitted all enabled statuses may match."`
}

func RegisterSearchRules(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[SearchRulesParams]()
	searchRulesTool := &mcp.Tool{
		Name:        "search_rules",
		Description: "Search for rules in the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, searchRulesTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[SearchRulesParams]) (*mcp.CallToolResultFor[any], error) {
		query := params.Arguments.Query

		opts := []any{}
		if params.Arguments.Anchoring != "" {
			opts = append(opts, opt.Anchoring(params.Arguments.Anchoring))
		}
		if params.Arguments.Context != "" {
			opts = append(opts, opt.RuleContexts(params.Arguments.Context))
		}
		if params.Arguments.Enabled != nil {
			opts = append(opts, opt.EnableRules(*params.Arguments.Enabled))
		}

		resp, err := index.SearchRules(query, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not search rules: %w", err)
		}

		return mcputil.JSONToolResult("rules", resp)
	})
}
