package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetReachability registers the get_reachability tool with the MCP server.
func RegisterGetReachability(mcps *server.MCPServer) {
	getReachabilityTool := mcp.NewTool(
		"monitoring_get_reachability",
		mcp.WithDescription("Test whether clusters are reachable or not"),
		mcp.WithString(
			"clusters",
			mcp.Description("Subset of clusters, separated by commas (e.g., c1-de,c2-de,c3-de)"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getReachabilityTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		clusters, _ := req.Params.Arguments["clusters"].(string)
		if clusters == "" {
			return nil, fmt.Errorf("clusters parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://status.algolia.com/1/reachability/%s/probes", clusters)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
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

		return mcputil.JSONToolResult("Reachability", result)
	})
}
