package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetTopSearches registers the get_top_searches tool with the MCP server.
func RegisterGetTopSearches(mcps *server.MCPServer) {
	getTopSearchesTool := mcp.NewTool(
		"analytics_get_top_searches",
		mcp.WithDescription("Retrieve the most popular searches for an index"),
		mcp.WithString(
			"index",
			mcp.Description("Index name"),
			mcp.Required(),
		),
		mcp.WithBoolean(
			"clickAnalytics",
			mcp.Description("Whether to include metrics related to click and conversion events in the response"),
		),
		mcp.WithBoolean(
			"revenueAnalytics",
			mcp.Description("Whether to include metrics related to revenue events in the response"),
		),
		mcp.WithString(
			"startDate",
			mcp.Description("Start date of the period to analyze, in YYYY-MM-DD format"),
		),
		mcp.WithString(
			"endDate",
			mcp.Description("End date of the period to analyze, in YYYY-MM-DD format"),
		),
		mcp.WithString(
			"orderBy",
			mcp.Description("Attribute by which to order the response items (searchCount, clickThroughRate, conversionRate, averageClickPosition)"),
		),
		mcp.WithString(
			"direction",
			mcp.Description("Sorting direction of the results: asc or desc"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Number of items to return (max 1000)"),
		),
		mcp.WithNumber(
			"offset",
			mcp.Description("Position of the first item to return"),
		),
		mcp.WithString(
			"tags",
			mcp.Description("Tags by which to segment the analytics"),
		),
	)

	mcps.AddTool(getTopSearchesTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		index, _ := req.Params.Arguments["index"].(string)
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

		if clickAnalytics, ok := req.Params.Arguments["clickAnalytics"].(bool); ok && clickAnalytics {
			q.Add("clickAnalytics", "true")
		}

		if revenueAnalytics, ok := req.Params.Arguments["revenueAnalytics"].(bool); ok && revenueAnalytics {
			q.Add("revenueAnalytics", "true")
		}

		if startDate, ok := req.Params.Arguments["startDate"].(string); ok && startDate != "" {
			q.Add("startDate", startDate)
		}

		if endDate, ok := req.Params.Arguments["endDate"].(string); ok && endDate != "" {
			q.Add("endDate", endDate)
		}

		if orderBy, ok := req.Params.Arguments["orderBy"].(string); ok && orderBy != "" {
			q.Add("orderBy", orderBy)
		}

		if direction, ok := req.Params.Arguments["direction"].(string); ok && direction != "" {
			q.Add("direction", direction)
		}

		if limit, ok := req.Params.Arguments["limit"].(float64); ok {
			q.Add("limit", strconv.FormatInt(int64(limit), 10))
		}

		if offset, ok := req.Params.Arguments["offset"].(float64); ok {
			q.Add("offset", strconv.FormatInt(int64(offset), 10))
		}

		if tags, ok := req.Params.Arguments["tags"].(string); ok && tags != "" {
			q.Add("tags", tags)
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

		return mcputil.JSONToolResult("Top Searches", result)
	})
}
