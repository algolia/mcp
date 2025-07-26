package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetIncidentsParams defines the parameters for retrieving incidents.
type GetIncidentsParams struct{}

// RegisterGetIncidents registers the get_incidents tool with the MCP server.
func RegisterGetIncidents(s *mcp.Server) {
	schema, _ := jsonschema.For[GetIncidentsParams]()
	getIncidentsTool := &mcp.Tool{
		Name:        "monitoring_get_incidents",
		Description: "Retrieves known incidents for all clusters",
		InputSchema: schema,
	}

	mcp.AddTool(s, getIncidentsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetIncidentsParams]) (*mcp.CallToolResultFor[any], error) {
		// Create HTTP client and request
		client := &http.Client{}
		url := "https://status.algolia.com/1/incidents"
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
					Text: "Incidents: " + fmt.Sprintf("%v", result),
				},
			},
		}, nil
	})
}
