package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetNoResultsRateParams defines the parameters for retrieving no results rate.
type GetNoResultsRateParams struct {
	Index     string `json:"index" jsonschema:"Index name"`
	StartDate string `json:"startDate,omitempty" jsonschema:"Start date of the period to analyze, in YYYY-MM-DD format"`
	EndDate   string `json:"endDate,omitempty" jsonschema:"End date of the period to analyze, in YYYY-MM-DD format"`
	Tags      string `json:"tags,omitempty" jsonschema:"Tags by which to segment the analytics"`
}

// RegisterGetNoResultsRate registers the get_no_results_rate tool with the MCP server.
func RegisterGetNoResultsRate(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetNoResultsRateParams]()
	getNoResultsRateTool := &mcp.Tool{
		Name:        "analytics_get_no_results_rate",
		Description: "Retrieve the fraction of searches that didn't return any results within a time range, including a daily breakdown",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getNoResultsRateTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetNoResultsRateParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		index := params.Arguments.Index
		if index == "" {
			return nil, fmt.Errorf("index parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://analytics.algolia.com/2/searches/noResultRate"
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("x-algolia-application-id", appID)
		httpReq.Header.Set("x-algolia-api-key", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		// Add query parameters
		q := httpReq.URL.Query()
		q.Add("index", index)

		if params.Arguments.StartDate != "" {
			q.Add("startDate", params.Arguments.StartDate)
		}

		if params.Arguments.EndDate != "" {
			q.Add("endDate", params.Arguments.EndDate)
		}

		if params.Arguments.Tags != "" {
			q.Add("tags", params.Arguments.Tags)
		}

		httpReq.URL.RawQuery = q.Encode()

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
					Text: "No Results Rate: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
