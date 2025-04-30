package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterListABTests registers the list_abtests tool with the MCP server.
func RegisterListABTests(mcps *server.MCPServer) {
	listABTestsTool := mcp.NewTool(
		"abtesting_list_abtests",
		mcp.WithDescription("List all A/B tests configured for this application"),
		mcp.WithNumber(
			"offset",
			mcp.Description("Position of the first item to return"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Number of items to return"),
		),
		mcp.WithString(
			"indexPrefix",
			mcp.Description("Index name prefix. Only A/B tests for indices starting with this string are included in the response"),
		),
		mcp.WithString(
			"indexSuffix",
			mcp.Description("Index name suffix. Only A/B tests for indices ending with this string are included in the response"),
		),
	)

	mcps.AddTool(listABTestsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Create Algolia Analytics client
		client := analytics.NewClient(appID, apiKey)

		// Prepare options
		opts := []interface{}{}

		if offset, ok := req.Params.Arguments["offset"].(float64); ok {
			opts = append(opts, opt.Offset(int(offset)))
		}

		if limit, ok := req.Params.Arguments["limit"].(float64); ok {
			opts = append(opts, opt.Limit(int(limit)))
		}

		if indexPrefix, ok := req.Params.Arguments["indexPrefix"].(string); ok && indexPrefix != "" {
			opts = append(opts, opt.IndexPrefix(indexPrefix))
		}

		if indexSuffix, ok := req.Params.Arguments["indexSuffix"].(string); ok && indexSuffix != "" {
			opts = append(opts, opt.IndexSuffix(indexSuffix))
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

		return mcputil.JSONToolResult("AB Tests", result)
	})
}
