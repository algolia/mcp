package collections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCollectionParams defines the parameters for retrieving a collection.
type GetCollectionParams struct {
	ID string `json:"id" jsonschema:"Collection ID"`
}

// RegisterGetCollection registers the get_collection tool with the MCP server.
func RegisterGetCollection(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetCollectionParams]()
	getCollectionTool := &mcp.Tool{
		Name:        "collections_get_collection",
		Description: "Retrieve a collection by ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getCollectionTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetCollectionParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		id := params.Arguments.ID
		if id == "" {
			return nil, fmt.Errorf("id parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://experiences.algolia.com/1/collections/%s", id)
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
					Text: fmt.Sprintf("Collection: %v", result),
				},
			},
		}, nil
	})
}
