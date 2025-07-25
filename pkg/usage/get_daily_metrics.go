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

// GetDailyMetricsParams defines the parameters for retrieving daily metrics.
type GetDailyMetricsParams struct {
	Applications string `json:"applications" jsonschema:"Comma-separated list of Algolia Application IDs"`
	StartDate    string `json:"startDate" jsonschema:"The start date of the period for which the metrics should be returned (YYYY-MM-DD)"`
	EndDate      string `json:"endDate,omitempty" jsonschema:"The end date (included) of the period for which the metrics should be returned (YYYY-MM-DD)"`
	MetricNames  string `json:"metricNames" jsonschema:"Comma-separated list of metric names to retrieve"`
}

// RegisterGetDailyMetrics registers the get_daily_metrics tool with the MCP server.
func RegisterGetDailyMetrics(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetDailyMetricsParams]()
	getDailyMetricsTool := &mcp.Tool{
		Name:        "usage_get_daily_metrics",
		Description: "Returns a list of billing metrics per day for the specified applications",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getDailyMetricsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetDailyMetricsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		applicationsStr := params.Arguments.Applications
		if applicationsStr == "" {
			return nil, fmt.Errorf("applications parameter is required")
		}

		startDate := params.Arguments.StartDate
		if startDate == "" {
			return nil, fmt.Errorf("startDate parameter is required")
		}

		metricNamesStr := params.Arguments.MetricNames
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
		urlParams := url.Values{}
		for _, app := range applications {
			urlParams.Add("application", app)
		}
		urlParams.Add("startDate", startDate)
		if params.Arguments.EndDate != "" {
			urlParams.Add("endDate", params.Arguments.EndDate)
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
					Text: "Daily Metrics: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
