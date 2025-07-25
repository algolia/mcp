package querysuggestions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteConfigParams defines the parameters for deleting a Query Suggestions configuration.
type DeleteConfigParams struct {
	Region    string `json:"region" jsonschema:"Analytics region (us or eu)"`
	IndexName string `json:"indexName" jsonschema:"Query Suggestions index name"`
}

// RegisterDeleteConfig registers the delete_query_suggestions_config tool with the MCP server.
func RegisterDeleteConfig(mcps *mcp.Server) {
	schema, _ := jsonschema.For[DeleteConfigParams]()
	deleteConfigTool := &mcp.Tool{
		Name:        "query_suggestions_delete_config",
		Description: "Deletes a Query Suggestions configuration",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteConfigTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteConfigParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for deleting configurations
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		region := params.Arguments.Region
		if region == "" {
			return nil, fmt.Errorf("region parameter is required")
		}

		indexName := params.Arguments.IndexName
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		// Validate region
		if region != "us" && region != "eu" {
			return nil, fmt.Errorf("region must be 'us' or 'eu'")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://query-suggestions.%s.algolia.com/1/configs/%s", region, indexName)
		httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
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
					Text: fmt.Sprintf("Query Suggestions Configuration Deleted: %v", result),
				},
			},
		}, nil
	})
}
