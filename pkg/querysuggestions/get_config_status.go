package querysuggestions

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

// RegisterGetConfigStatus registers the get_query_suggestions_config_status tool with the MCP server.
func RegisterGetConfigStatus(mcps *server.MCPServer) {
	getConfigStatusTool := mcp.NewTool(
		"query_suggestions_get_config_status",
		mcp.WithDescription("Reports the status of a Query Suggestions index"),
		mcp.WithString(
			"region",
			mcp.Description("Analytics region (us or eu)"),
			mcp.Required(),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("Query Suggestions index name"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getConfigStatusTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		region, _ := req.Params.Arguments["region"].(string)
		if region == "" {
			return nil, fmt.Errorf("region parameter is required")
		}

		indexName, _ := req.Params.Arguments["indexName"].(string)
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		// Validate region
		if region != "us" && region != "eu" {
			return nil, fmt.Errorf("region must be 'us' or 'eu'")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://query-suggestions.%s.algolia.com/1/configs/%s/status", region, indexName)
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

		return mcputil.JSONToolResult("Query Suggestions Configuration Status", result)
	})
}
