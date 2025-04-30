package monitoring

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

// RegisterGetMetrics registers the get_metrics tool with the MCP server.
func RegisterGetMetrics(mcps *server.MCPServer) {
	getMetricsTool := mcp.NewTool(
		"monitoring_get_metrics",
		mcp.WithDescription("Retrieves metrics related to your Algolia infrastructure, aggregated over a selected time window"),
		mcp.WithString(
			"metric",
			mcp.Description("Metric to report (avg_build_time, ssd_usage, ram_search_usage, ram_indexing_usage, cpu_usage, or * for all)"),
			mcp.Required(),
		),
		mcp.WithString(
			"period",
			mcp.Description("Period over which to aggregate the metrics (minute, hour, day, week, month)"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getMetricsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		metric, _ := req.Params.Arguments["metric"].(string)
		if metric == "" {
			return nil, fmt.Errorf("metric parameter is required")
		}

		period, _ := req.Params.Arguments["period"].(string)
		if period == "" {
			return nil, fmt.Errorf("period parameter is required")
		}

		// Validate metric
		validMetrics := map[string]bool{
			"avg_build_time":     true,
			"ssd_usage":          true,
			"ram_search_usage":   true,
			"ram_indexing_usage": true,
			"cpu_usage":          true,
			"*":                  true,
		}
		if !validMetrics[metric] {
			return nil, fmt.Errorf("invalid metric: %s", metric)
		}

		// Validate period
		validPeriods := map[string]bool{
			"minute": true,
			"hour":   true,
			"day":    true,
			"week":   true,
			"month":  true,
		}
		if !validPeriods[period] {
			return nil, fmt.Errorf("invalid period: %s", period)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://status.algolia.com/1/infrastructure/%s/period/%s", metric, period)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("X-ALGOLIA-APPLICATION-ID", appID)
		httpReq.Header.Set("X-ALGOLIA-API-KEY", apiKey)
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

		return mcputil.JSONToolResult("Infrastructure Metrics", result)
	})
}
