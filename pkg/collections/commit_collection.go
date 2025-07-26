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

// CommitCollectionParams defines the parameters for committing a collection.
type CommitCollectionParams struct {
	ID string `json:"id" jsonschema:"Collection ID"`
}

// RegisterCommitCollection registers the commit_collection tool with the MCP server.
func RegisterCommitCollection(mcps *mcp.Server) {
	schema, _ := jsonschema.For[CommitCollectionParams]()
	commitCollectionTool := &mcp.Tool{
		Name:        "collections_commit_collection",
		Description: "Evaluates the changes on a collection and replicates them to the index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, commitCollectionTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[CommitCollectionParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for committing collections
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		id := params.Arguments.ID
		if id == "" {
			return nil, fmt.Errorf("id parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://experiences.algolia.com/1/collections/%s/commit", id)
		httpReq, err := http.NewRequest(http.MethodPost, url, nil)
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
		if resp.StatusCode != http.StatusAccepted {
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
					Text: fmt.Sprintf("Collection Commit Started: %v", result),
				},
			},
		}, nil
	})
}
