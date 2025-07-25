package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetHourlyMetricsParams defines the parameters for retrieving hourly metrics.
type GetHourlyMetricsParams struct {
	Application string `json:"application" jsonschema:"Algolia Application ID"`
	StartTime   string `json:"startTime" jsonschema:"The start time of the period for which the metrics should be returned (ISO 8601 format)"`
	EndTime     string `json:"endTime,omitempty" jsonschema:"The end time (included) of the period for which the metrics should be returned (ISO 8601 format)"`
	MetricNames string `json:"metricNames" jsonschema:"Comma-separated list of metric names to retrieve"`
}

// RegisterGetHourlyMetrics registers the get_hourly_metrics tool with the MCP server.
func RegisterGetHourlyMetrics(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetHourlyMetricsParams]()
	getHourlyMetricsTool := &mcp.Tool{
		Name:        "usage_get_hourly_metrics",
		Description: "Returns a list of billing metrics per hour for the specified application",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getHourlyMetricsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetHourlyMetricsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		application := params.Arguments.Application
		if application == "" {
			return nil, fmt.Errorf("application parameter is required")
		}

		startTime := params.Arguments.StartTime
		if startTime == "" {
			return nil, fmt.Errorf("startTime parameter is required")
		}

		metricNamesStr := params.Arguments.MetricNames
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
		urlParams := url.Values{}
		urlParams.Add("application", application)
		urlParams.Add("startTime", startTime)
		if params.Arguments.EndTime != "" {
			urlParams.Add("endTime", params.Arguments.EndTime)
		}
		for _, name := range metricNames {
			urlParams.Add("name", name)
		}
		url := fmt.Sprintf("%s?%s", baseURL, urlParams.Encode())

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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Hourly Metrics: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
