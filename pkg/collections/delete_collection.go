package collections

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

// RegisterDeleteCollection registers the delete_collection tool with the MCP server.
func RegisterDeleteCollection(mcps *server.MCPServer) {
	deleteCollectionTool := mcp.NewTool(
		"collections_delete_collection",
		mcp.WithDescription("Soft deletes a collection by setting 'deleted' to true"),
		mcp.WithString(
			"id",
			mcp.Description("Collection ID"),
			mcp.Required(),
		),
	)

	mcps.AddTool(deleteCollectionTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for deleting collections
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		id, _ := req.Params.Arguments["id"].(string)
		if id == "" {
			return nil, fmt.Errorf("id parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://experiences.algolia.com/1/collections/%s", id)
		httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
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

		return mcputil.JSONToolResult("Collection Deleted", map[string]any{
			"id":      id,
			"deleted": true,
		})
	})
}
