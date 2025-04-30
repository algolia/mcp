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

// RegisterStopABTest registers the stop_abtest tool with the MCP server.
func RegisterStopABTest(mcps *server.MCPServer) {
	stopABTestTool := mcp.NewTool(
		"abtesting_stop_abtest",
		mcp.WithDescription("Stop an A/B test by its ID. You can't restart stopped A/B tests."),
		mcp.WithNumber(
			"id",
			mcp.Description("Unique A/B test identifier"),
			mcp.Required(),
		),
	)

	mcps.AddTool(stopABTestTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for stopping AB tests
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Get the AB Test ID from the request
		idFloat, ok := req.Params.Arguments["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid AB test ID")
		}
		id := int(idFloat)

		// Create Algolia Analytics client
		client := analytics.NewClient(appID, apiKey)

		// Stop AB test
		res, err := client.StopABTest(id)
		if err != nil {
			return nil, fmt.Errorf("failed to stop AB test: %w", err)
		}

		// Convert to map for consistent response format
		result := map[string]interface{}{
			"taskID": res.TaskID,
			"index":  res.Index,
		}

		return mcputil.JSONToolResult(fmt.Sprintf("AB Test %d Stopped", id), result)
	})
}
