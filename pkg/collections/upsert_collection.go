package collections

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterUpsertCollection registers the upsert_collection tool with the MCP server.
func RegisterUpsertCollection(mcps *server.MCPServer) {
	upsertCollectionTool := mcp.NewTool(
		"collections_upsert_collection",
		mcp.WithDescription("Upserts a collection"),
		mcp.WithString(
			"id",
			mcp.Description("Collection ID (optional for new collections)"),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("Name of the index"),
			mcp.Required(),
		),
		mcp.WithString(
			"name",
			mcp.Description("Collection name"),
			mcp.Required(),
		),
		mcp.WithString(
			"description",
			mcp.Description("Collection description"),
		),
		mcp.WithString(
			"add",
			mcp.Description("JSON array of objectIDs to add to the collection"),
		),
		mcp.WithString(
			"remove",
			mcp.Description("JSON array of objectIDs to remove from the collection"),
		),
		mcp.WithString(
			"conditions",
			mcp.Description("JSON object with conditions to filter records"),
		),
	)

	mcps.AddTool(upsertCollectionTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for creating/updating collections
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract required parameters
		indexName, _ := req.Params.Arguments["indexName"].(string)
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		name, _ := req.Params.Arguments["name"].(string)
		if name == "" {
			return nil, fmt.Errorf("name parameter is required")
		}

		// Prepare request body
		requestBody := map[string]any{
			"indexName": indexName,
			"name":      name,
		}

		// Add optional parameters if provided
		if id, ok := req.Params.Arguments["id"].(string); ok && id != "" {
			requestBody["id"] = id
		}

		if description, ok := req.Params.Arguments["description"].(string); ok && description != "" {
			requestBody["description"] = description
		}

		// Parse and add 'add' array if provided
		if addJSON, ok := req.Params.Arguments["add"].(string); ok && addJSON != "" {
			var add []string
			if err := json.Unmarshal([]byte(addJSON), &add); err != nil {
				return nil, fmt.Errorf("invalid add JSON: %w", err)
			}
			requestBody["add"] = add
		}

		// Parse and add 'remove' array if provided
		if removeJSON, ok := req.Params.Arguments["remove"].(string); ok && removeJSON != "" {
			var remove []string
			if err := json.Unmarshal([]byte(removeJSON), &remove); err != nil {
				return nil, fmt.Errorf("invalid remove JSON: %w", err)
			}
			requestBody["remove"] = remove
		}

		// Parse and add 'conditions' object if provided
		if conditionsJSON, ok := req.Params.Arguments["conditions"].(string); ok && conditionsJSON != "" {
			var conditions map[string]any
			if err := json.Unmarshal([]byte(conditionsJSON), &conditions); err != nil {
				return nil, fmt.Errorf("invalid conditions JSON: %w", err)
			}
			requestBody["conditions"] = conditions
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://experiences.algolia.com/1/collections"
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
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
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

		return mcputil.JSONToolResult("Collection Upserted", result)
	})
}
