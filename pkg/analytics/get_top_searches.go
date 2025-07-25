package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetTopSearchesParams defines the parameters for retrieving top searches.
type GetTopSearchesParams struct {
	Index            string  `json:"index" jsonschema:"Index name"`
	ClickAnalytics   *bool   `json:"clickAnalytics,omitempty" jsonschema:"Whether to include metrics related to click and conversion events in the response"`
	RevenueAnalytics *bool   `json:"revenueAnalytics,omitempty" jsonschema:"Whether to include metrics related to revenue events in the response"`
	StartDate        string  `json:"startDate,omitempty" jsonschema:"Start date of the period to analyze, in YYYY-MM-DD format"`
	EndDate          string  `json:"endDate,omitempty" jsonschema:"End date of the period to analyze, in YYYY-MM-DD format"`
	OrderBy          string  `json:"orderBy,omitempty" jsonschema:"Attribute by which to order the response items (searchCount, clickThroughRate, conversionRate, averageClickPosition)"`
	Direction        string  `json:"direction,omitempty" jsonschema:"Sorting direction of the results: asc or desc"`
	Limit            *int    `json:"limit,omitempty" jsonschema:"Number of items to return (max 1000)"`
	Offset           *int    `json:"offset,omitempty" jsonschema:"Position of the first item to return"`
	Tags             string  `json:"tags,omitempty" jsonschema:"Tags by which to segment the analytics"`
}

// RegisterGetTopSearches registers the get_top_searches tool with the MCP server.
func RegisterGetTopSearches(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetTopSearchesParams]()
	getTopSearchesTool := &mcp.Tool{
		Name:        "analytics_get_top_searches",
		Description: "Retrieve the most popular searches for an index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getTopSearchesTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetTopSearchesParams]) (*mcp.CallToolResultFor[any], error) {
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
		url := "https://analytics.algolia.com/2/searches"
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

		if params.Arguments.ClickAnalytics != nil && *params.Arguments.ClickAnalytics {
			q.Add("clickAnalytics", "true")
		}

		if params.Arguments.RevenueAnalytics != nil && *params.Arguments.RevenueAnalytics {
			q.Add("revenueAnalytics", "true")
		}

		if params.Arguments.StartDate != "" {
			q.Add("startDate", params.Arguments.StartDate)
		}

		if params.Arguments.EndDate != "" {
			q.Add("endDate", params.Arguments.EndDate)
		}

		if params.Arguments.OrderBy != "" {
			q.Add("orderBy", params.Arguments.OrderBy)
		}

		if params.Arguments.Direction != "" {
			q.Add("direction", params.Arguments.Direction)
		}

		if params.Arguments.Limit != nil {
			q.Add("limit", strconv.Itoa(*params.Arguments.Limit))
		}

		if params.Arguments.Offset != nil {
			q.Add("offset", strconv.Itoa(*params.Arguments.Offset))
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
					Text: "Top Searches: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
