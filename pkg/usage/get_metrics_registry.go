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

// RegisterGetMetricsRegistry registers the get_metrics_registry tool with the MCP server.
func RegisterGetMetricsRegistry(mcps *server.MCPServer) {
	getMetricsRegistryTool := mcp.NewTool(
		"usage_get_metrics_registry",
		mcp.WithDescription("Returns the list of available metrics"),
		mcp.WithString(
			"applications",
			mcp.Description("Comma-separated list of Algolia Application IDs"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getMetricsRegistryTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Split applications string into array
		applications := strings.Split(applicationsStr, ",")
		for i, app := range applications {
			applications[i] = strings.TrimSpace(app)
		}

		// Create HTTP client and request
		client := &http.Client{}
		baseURL := "https://usage.algolia.com/2/metrics/registry"

		// Add query parameters
		params := url.Values{}
		for _, app := range applications {
			params.Add("application", app)
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

		return mcputil.JSONToolResult("Metrics Registry", result)
	})
}
