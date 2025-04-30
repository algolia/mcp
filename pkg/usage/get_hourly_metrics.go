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

// RegisterGetHourlyMetrics registers the get_hourly_metrics tool with the MCP server.
func RegisterGetHourlyMetrics(mcps *server.MCPServer) {
	getHourlyMetricsTool := mcp.NewTool(
		"usage_get_hourly_metrics",
		mcp.WithDescription("Returns a list of billing metrics per hour for the specified application"),
		mcp.WithString(
			"application",
			mcp.Description("Algolia Application ID"),
			mcp.Required(),
		),
		mcp.WithString(
			"startTime",
			mcp.Description("The start time of the period for which the metrics should be returned (ISO 8601 format)"),
			mcp.Required(),
		),
		mcp.WithString(
			"endTime",
			mcp.Description("The end time (included) of the period for which the metrics should be returned (ISO 8601 format)"),
		),
		mcp.WithString(
			"metricNames",
			mcp.Description("Comma-separated list of metric names to retrieve"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getHourlyMetricsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		application, _ := req.Params.Arguments["application"].(string)
		if application == "" {
			return nil, fmt.Errorf("application parameter is required")
		}

		startTime, _ := req.Params.Arguments["startTime"].(string)
		if startTime == "" {
			return nil, fmt.Errorf("startTime parameter is required")
		}

		metricNamesStr, _ := req.Params.Arguments["metricNames"].(string)
		if metricNamesStr == "" {
			return nil, fmt.Errorf("metricNames parameter is required")
		}

		// Split metric names string into array
		metricNames := strings.Split(metricNamesStr, ",")
		for i, name := range metricNames {
			metricNames[i] = strings.TrimSpace(name)
		}

		// Create HTTP client and request
		client := &http.Client{}
		baseURL := "https://usage.algolia.com/2/metrics/hourly"

		// Add query parameters
		params := url.Values{}
		params.Add("application", application)
		params.Add("startTime", startTime)
		if endTime, ok := req.Params.Arguments["endTime"].(string); ok && endTime != "" {
			params.Add("endTime", endTime)
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

		return mcputil.JSONToolResult("Hourly Metrics", result)
	})
}
