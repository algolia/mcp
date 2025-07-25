package collections

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// UpsertCollectionParams defines the parameters for upserting a collection.
type UpsertCollectionParams struct {
	ID          *string `json:"id,omitempty" jsonschema:"Collection ID (optional for new collections)"`
	IndexName   string  `json:"indexName" jsonschema:"Name of the index"`
	Name        string  `json:"name" jsonschema:"Collection name"`
	Description *string `json:"description,omitempty" jsonschema:"Collection description"`
	Add         *string `json:"add,omitempty" jsonschema:"JSON array of objectIDs to add to the collection"`
	Remove      *string `json:"remove,omitempty" jsonschema:"JSON array of objectIDs to remove from the collection"`
	Conditions  *string `json:"conditions,omitempty" jsonschema:"JSON object with conditions to filter records"`
}

// RegisterUpsertCollection registers the upsert_collection tool with the MCP server.
func RegisterUpsertCollection(mcps *mcp.Server) {
	schema, _ := jsonschema.For[UpsertCollectionParams]()
	upsertCollectionTool := &mcp.Tool{
		Name:        "collections_upsert_collection",
		Description: "Upserts a collection",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, upsertCollectionTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[UpsertCollectionParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for creating/updating collections
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract required parameters
		indexName := params.Arguments.IndexName
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		name := params.Arguments.Name
		if name == "" {
			return nil, fmt.Errorf("name parameter is required")
		}

		// Prepare request body
		requestBody := map[string]any{
			"indexName": indexName,
			"name":      name,
		}

		// Add optional parameters if provided
		if params.Arguments.ID != nil && *params.Arguments.ID != "" {
			requestBody["id"] = *params.Arguments.ID
		}

		if params.Arguments.Description != nil && *params.Arguments.Description != "" {
			requestBody["description"] = *params.Arguments.Description
		}

		// Parse and add 'add' array if provided
		if params.Arguments.Add != nil && *params.Arguments.Add != "" {
			var add []string
			if err := json.Unmarshal([]byte(*params.Arguments.Add), &add); err != nil {
				return nil, fmt.Errorf("invalid add JSON: %w", err)
			}
			requestBody["add"] = add
		}

		// Parse and add 'remove' array if provided
		if params.Arguments.Remove != nil && *params.Arguments.Remove != "" {
			var remove []string
			if err := json.Unmarshal([]byte(*params.Arguments.Remove), &remove); err != nil {
				return nil, fmt.Errorf("invalid remove JSON: %w", err)
			}
			requestBody["remove"] = remove
		}

		// Parse and add 'conditions' object if provided
		if params.Arguments.Conditions != nil && *params.Arguments.Conditions != "" {
			var conditions map[string]any
			if err := json.Unmarshal([]byte(*params.Arguments.Conditions), &conditions); err != nil {
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Collection Upserted: %v", result),
				},
			},
		}, nil
	})
}
