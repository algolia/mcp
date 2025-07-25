package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteABTestParams defines the parameters for deleting an A/B test.
type DeleteABTestParams struct {
	ID float64 `json:"id" jsonschema:"Unique A/B test identifier"`
}

// RegisterDeleteABTest registers the delete_abtest tool with the MCP server.
func RegisterDeleteABTest(mcps *mcp.Server) {
	schema, _ := jsonschema.For[DeleteABTestParams]()
	deleteABTestTool := &mcp.Tool{
		Name:        "abtesting_delete_abtest",
		Description: "Delete an A/B test by its ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteABTestTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteABTestParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for deleting AB tests
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Get the AB Test ID from the request
		id := int(params.Arguments.ID)

		// Create Algolia Analytics client
		client := analytics.NewClient(appID, apiKey)

		// Delete AB test
		res, err := client.DeleteABTest(id)
		if err != nil {
			return nil, fmt.Errorf("failed to delete AB test: %w", err)
		}

		// Convert to map for consistent response format
		result := map[string]interface{}{
			"taskID": res.TaskID,
			"index":  res.Index,
		}

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("AB Test %d Deleted: %v", id, result),
				},
			},
		}, nil
	})
}
