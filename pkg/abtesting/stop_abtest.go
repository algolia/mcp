package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StopABTestParams defines the parameters for stopping an A/B test.
type StopABTestParams struct {
	ID float64 `json:"id" jsonschema:"Unique A/B test identifier"`
}

// RegisterStopABTest registers the stop_abtest tool with the MCP server.
func RegisterStopABTest(mcps *mcp.Server) {
	schema, _ := jsonschema.For[StopABTestParams]()
	stopABTestTool := &mcp.Tool{
		Name:        "abtesting_stop_abtest",
		Description: "Stop an A/B test by its ID. You can't restart stopped A/B tests.",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, stopABTestTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[StopABTestParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for stopping AB tests
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Get the AB Test ID from the request
		id := int(params.Arguments.ID)

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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("AB Test %d Stopped: %v", id, result),
				},
			},
		}, nil
	})
}
