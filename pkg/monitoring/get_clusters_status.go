package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetClustersStatusParams defines the parameters for retrieving clusters status.
type GetClustersStatusParams struct{}

// RegisterGetClustersStatus registers the get_clusters_status tool with the MCP server.
func RegisterGetClustersStatus(s *mcp.Server) {
	schema, _ := jsonschema.For[GetClustersStatusParams]()
	getClustersStatusTool := &mcp.Tool{
		Name:        "monitoring_get_clusters_status",
		Description: "Retrieves the status of all Algolia clusters and instances",
		InputSchema: schema,
	}

	mcp.AddTool(s, getClustersStatusTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetClustersStatusParams]) (*mcp.CallToolResultFor[any], error) {
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Clusters Status: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
