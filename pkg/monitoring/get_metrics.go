package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetMetricsParams defines the parameters for retrieving metrics.
type GetMetricsParams struct {
	Metric string `json:"metric" jsonschema:"Metric to report (avg_build_time, ssd_usage, ram_search_usage, ram_indexing_usage, cpu_usage, or * for all)"`
	Period string `json:"period" jsonschema:"Period over which to aggregate the metrics (minute, hour, day, week, month)"`
}

// RegisterGetMetrics registers the get_metrics tool with the MCP server.
func RegisterGetMetrics(s *mcp.Server) {
	schema, _ := jsonschema.For[GetMetricsParams]()
	getMetricsTool := &mcp.Tool{
		Name:        "monitoring_get_metrics",
		Description: "Retrieves metrics related to your Algolia infrastructure, aggregated over a selected time window",
		InputSchema: schema,
	}

	mcp.AddTool(s, getMetricsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetMetricsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		metric := params.Arguments.Metric
		if metric == "" {
			return nil, fmt.Errorf("metric parameter is required")
		}

		period := params.Arguments.Period
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Infrastructure Metrics: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
