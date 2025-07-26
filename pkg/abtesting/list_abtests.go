package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListABTestsParams defines the parameters for listing A/B tests.
type ListABTestsParams struct {
	Offset      *float64 `json:"offset,omitempty" jsonschema:"Position of the first item to return"`
	Limit       *float64 `json:"limit,omitempty" jsonschema:"Number of items to return"`
	IndexPrefix *string  `json:"indexPrefix,omitempty" jsonschema:"Index name prefix. Only A/B tests for indices starting with this string are included in the response"`
	IndexSuffix *string  `json:"indexSuffix,omitempty" jsonschema:"Index name suffix. Only A/B tests for indices ending with this string are included in the response"`
}

// RegisterListABTests registers the list_abtests tool with the MCP server.
func RegisterListABTests(mcps *mcp.Server) {
	schema, _ := jsonschema.For[ListABTestsParams]()
	listABTestsTool := &mcp.Tool{
		Name:        "abtesting_list_abtests",
		Description: "List all A/B tests configured for this application",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, listABTestsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListABTestsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Create Algolia Analytics client
		client := analytics.NewClient(appID, apiKey)

		// Prepare options
		opts := []interface{}{}

		if params.Arguments.Offset != nil {
			opts = append(opts, opt.Offset(int(*params.Arguments.Offset)))
		}

		if params.Arguments.Limit != nil {
			opts = append(opts, opt.Limit(int(*params.Arguments.Limit)))
		}

		if params.Arguments.IndexPrefix != nil && *params.Arguments.IndexPrefix != "" {
			opts = append(opts, opt.IndexPrefix(*params.Arguments.IndexPrefix))
		}

		if params.Arguments.IndexSuffix != nil && *params.Arguments.IndexSuffix != "" {
			opts = append(opts, opt.IndexSuffix(*params.Arguments.IndexSuffix))
		}

		// Get AB tests
		res, err := client.GetABTests(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to get AB tests: %w", err)
		}

		// Convert to map for consistent response format
		result := map[string]interface{}{
			"count":   res.Count,
			"total":   res.Total,
			"abtests": res.ABTests,
		}

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("AB Tests: %v", result),
				},
			},
		}, nil
	})
}
