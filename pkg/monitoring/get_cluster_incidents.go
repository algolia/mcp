package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetClusterIncidentsParams defines the parameters for retrieving cluster incidents.
type GetClusterIncidentsParams struct {
	Clusters string `json:"clusters" jsonschema:"Subset of clusters, separated by commas (e.g., c1-de,c2-de,c3-de)"`
}

// RegisterGetClusterIncidents registers the get_cluster_incidents tool with the MCP server.
func RegisterGetClusterIncidents(s *mcp.Server) {
	schema, _ := jsonschema.For[GetClusterIncidentsParams]()
	getClusterIncidentsTool := &mcp.Tool{
		Name:        "monitoring_get_cluster_incidents",
		Description: "Retrieves known incidents for the selected clusters",
		InputSchema: schema,
	}

	mcp.AddTool(s, getClusterIncidentsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetClusterIncidentsParams]) (*mcp.CallToolResultFor[any], error) {
		// Extract parameters
		clusters := params.Arguments.Clusters
		if clusters == "" {
			return nil, fmt.Errorf("clusters parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://status.algolia.com/1/incidents/%s", clusters)
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
					Text: "Cluster Incidents: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
