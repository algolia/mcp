package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetClickThroughRate registers the get_click_through_rate tool with the MCP server
func RegisterGetClickThroughRate(mcps *server.MCPServer) {
	getClickThroughRateTool := mcp.NewTool(
		"analytics_get_click_through_rate",
		mcp.WithDescription("Retrieve the click-through rate (CTR) for all your searches with at least one click event, including a daily breakdown"),
		mcp.WithString(
			"index",
			mcp.Description("Index name"),
			mcp.Required(),
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
			"tags",
			mcp.Description("Tags by which to segment the analytics"),
		),
	)

	mcps.AddTool(getClickThroughRateTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		url := "https://analytics.algolia.com/2/clicks/clickThroughRate"
		httpReq, err := http.NewRequest("GET", url, nil)
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

		if startDate, ok := req.Params.Arguments["startDate"].(string); ok && startDate != "" {
			q.Add("startDate", startDate)
		}

		if endDate, ok := req.Params.Arguments["endDate"].(string); ok && endDate != "" {
			q.Add("endDate", endDate)
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

		return mcputil.JSONToolResult("Click Through Rate", result)
	})
}
