package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetABTest registers the get_abtest tool with the MCP server.
func RegisterGetABTest(mcps *server.MCPServer) {
	getABTestTool := mcp.NewTool(
		"abtesting_get_abtest",
		mcp.WithDescription("Retrieve the details for an A/B test by its ID"),
		mcp.WithNumber(
			"id",
			mcp.Description("Unique A/B test identifier"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getABTestTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Get the AB Test ID from the request
		idFloat, ok := req.Params.Arguments["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid AB test ID")
		}
		id := int(idFloat)

		// Create Algolia Analytics client
		client := analytics.NewClient(appID, apiKey)

		// Get AB test
		res, err := client.GetABTest(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get AB test: %w", err)
		}

		// Convert to map for consistent response format
		result := map[string]interface{}{
			"abTestID":               res.ABTestID,
			"clickSignificance":      res.ClickSignificance,
			"conversionSignificance": res.ConversionSignificance,
			"createdAt":              res.CreatedAt,
			"updatedAt":              res.UpdatedAt,
			"endAt":                  res.EndAt,
			"name":                   res.Name,
			"status":                 res.Status,
			"variants":               res.Variants,
		}

		return mcputil.JSONToolResult(fmt.Sprintf("AB Test %d", id), result)
	})
}
