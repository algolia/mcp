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

// ScheduleABTestParams defines the parameters for scheduling an A/B test.
type ScheduleABTestParams struct {
	Name        string `json:"name" jsonschema:"A/B test name"`
	ScheduledAt string `json:"scheduledAt" jsonschema:"Date and time when the A/B test is scheduled to start in RFC 3339 format (e.g. 2023-06-15T15:06:44.400601Z)"`
	EndAt       string `json:"endAt" jsonschema:"End date and time of the A/B test in RFC 3339 format (e.g. 2023-06-17T00:00:00Z)"`
	Variants    string `json:"variants" jsonschema:"A/B test variants as JSON array (exactly 2 variants required). Each variant must have 'index' and 'trafficPercentage' fields and may optionally have 'description' and 'customSearchParameters' fields."`
}

// RegisterScheduleABTest registers the schedule_abtest tool with the MCP server.
func RegisterScheduleABTest(mcps *mcp.Server) {
	schema, _ := jsonschema.For[ScheduleABTestParams]()
	scheduleABTestTool := &mcp.Tool{
		Name:        "abtesting_schedule_abtest",
		Description: "Schedule an A/B test to be started at a later time",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, scheduleABTestTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ScheduleABTestParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for scheduling AB tests
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		name := params.Arguments.Name
		scheduledAt := params.Arguments.ScheduledAt
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
			"name":        name,
			"scheduledAt": scheduledAt,
			"endAt":       endAt,
			"variants":    variants,
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://analytics.algolia.com/2/abtests/schedule"
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
					Text: fmt.Sprintf("AB Test Scheduled: %v", result),
				},
			},
		}, nil
	})
}
