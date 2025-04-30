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

// RegisterGetClustersStatus registers the get_clusters_status tool with the MCP server.
func RegisterGetClustersStatus(mcps *server.MCPServer) {
	getClustersStatusTool := mcp.NewTool(
		"monitoring_get_clusters_status",
		mcp.WithDescription("Retrieves the status of all Algolia clusters and instances"),
	)

	mcps.AddTool(getClustersStatusTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create HTTP client and request
		client := &http.Client{}
		url := "https://status.algolia.com/1/status"
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

		return mcputil.JSONToolResult("Clusters Status", result)
	})
}
