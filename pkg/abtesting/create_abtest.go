package abtesting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CreateABTestParams defines the parameters for creating an A/B test.
type CreateABTestParams struct {
	Name     string `json:"name" jsonschema:"A/B test name"`
	EndAt    string `json:"endAt" jsonschema:"End date and time of the A/B test in RFC 3339 format (e.g. 2023-06-17T00:00:00Z)"`
	Variants string `json:"variants" jsonschema:"A/B test variants as JSON array (exactly 2 variants required). Each variant must have 'index' and 'trafficPercentage' fields and may optionally have 'description' and 'customSearchParameters' fields."`
}

// RegisterCreateABTest registers the create_abtest tool with the MCP server.
func RegisterCreateABTest(mcps *mcp.Server) {
	schema, _ := jsonschema.For[CreateABTestParams]()
	createABTestTool := &mcp.Tool{
		Name:        "abtesting_create_abtest",
		Description: "Create a new A/B test",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, createABTestTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateABTestParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for creating AB tests
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		name := params.Arguments.Name
		endAt := params.Arguments.EndAt
		variantsJSON := params.Arguments.Variants

		// Parse variants JSON
		var variants []any
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			return nil, fmt.Errorf("invalid variants JSON: %w", err)
		}

		if len(variants) != 2 {
			return nil, fmt.Errorf("exactly 2 variants are required")
		}

		// Prepare request body
		requestBody := map[string]any{
			"name":     name,
			"endAt":    endAt,
			"variants": variants,
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://analytics.algolia.com/2/abtests"
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("x-algolia-application-id", appID)
		httpReq.Header.Set("x-algolia-api-key", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		// Check for error response
		if resp.StatusCode != http.StatusOK {
			var errResp map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
				return nil, fmt.Errorf("Algolia API error (status %d)", resp.StatusCode)
			}
			return nil, fmt.Errorf("Algolia API error: %v", errResp)
		}

		// Parse response
		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "AB Test Created: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
