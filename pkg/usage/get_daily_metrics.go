package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetDailyMetrics registers the get_daily_metrics tool with the MCP server.
func RegisterGetDailyMetrics(mcps *server.MCPServer) {
	getDailyMetricsTool := mcp.NewTool(
		"usage_get_daily_metrics",
		mcp.WithDescription("Returns a list of billing metrics per day for the specified applications"),
		mcp.WithString(
			"applications",
			mcp.Description("Comma-separated list of Algolia Application IDs"),
			mcp.Required(),
		),
		mcp.WithString(
			"startDate",
			mcp.Description("The start date of the period for which the metrics should be returned (YYYY-MM-DD)"),
			mcp.Required(),
		),
		mcp.WithString(
			"endDate",
			mcp.Description("The end date (included) of the period for which the metrics should be returned (YYYY-MM-DD)"),
		),
		mcp.WithString(
			"metricNames",
			mcp.Description("Comma-separated list of metric names to retrieve"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getDailyMetricsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		applicationsStr, _ := req.Params.Arguments["applications"].(string)
		if applicationsStr == "" {
			return nil, fmt.Errorf("applications parameter is required")
		}

		startDate, _ := req.Params.Arguments["startDate"].(string)
		if startDate == "" {
			return nil, fmt.Errorf("startDate parameter is required")
		}

		metricNamesStr, _ := req.Params.Arguments["metricNames"].(string)
		if metricNamesStr == "" {
			return nil, fmt.Errorf("metricNames parameter is required")
		}

		// Split applications string into array
		applications := strings.Split(applicationsStr, ",")
		for i, app := range applications {
			applications[i] = strings.TrimSpace(app)
		}

		// Split metric names string into array
		metricNames := strings.Split(metricNamesStr, ",")
		for i, name := range metricNames {
			metricNames[i] = strings.TrimSpace(name)
		}

		// Create HTTP client and request
		client := &http.Client{}
		baseURL := "https://usage.algolia.com/2/metrics/daily"

		// Add query parameters
		params := url.Values{}
		for _, app := range applications {
			params.Add("application", app)
		}
		params.Add("startDate", startDate)
		if endDate, ok := req.Params.Arguments["endDate"].(string); ok && endDate != "" {
			params.Add("endDate", endDate)
		}
		for _, name := range metricNames {
			params.Add("name", name)
		}
		url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
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

		return mcputil.JSONToolResult("Daily Metrics", result)
	})
}
