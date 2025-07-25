package abtesting

import (
	"context"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/analytics"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetABTestParams defines the parameters for retrieving an A/B test.
type GetABTestParams struct {
	ID float64 `json:"id" jsonschema:"Unique A/B test identifier"`
}

// RegisterGetABTest registers the get_abtest tool with the MCP server.
func RegisterGetABTest(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetABTestParams]()
	getABTestTool := &mcp.Tool{
		Name:        "abtesting_get_abtest",
		Description: "Retrieve the details for an A/B test by its ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getABTestTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetABTestParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Get the AB Test ID from the request
		id := int(params.Arguments.ID)

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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("AB Test %d: %v", id, result),
				},
			},
		}, nil
	})
}
