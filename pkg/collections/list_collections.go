package collections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListCollectionsParams defines the parameters for listing collections.
type ListCollectionsParams struct {
	IndexName string  `json:"indexName" jsonschema:"Name of the index"`
	Offset    *int    `json:"offset,omitempty" jsonschema:"Number of items to skip (default to 0)"`
	Limit     *int    `json:"limit,omitempty" jsonschema:"Number of items per fetch (defaults to 10)"`
	Query     *string `json:"query,omitempty" jsonschema:"Query to filter collections"`
}

// RegisterListCollections registers the list_collections tool with the MCP server.
func RegisterListCollections(mcps *mcp.Server) {
	schema, _ := jsonschema.For[ListCollectionsParams]()
	listCollectionsTool := &mcp.Tool{
		Name:        "collections_list_collections",
		Description: "Retrieve a list of all collections",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, listCollectionsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListCollectionsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		indexName := params.Arguments.IndexName
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://experiences.algolia.com/1/collections"
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("X-ALGOLIA-APPLICATION-ID", appID)
		httpReq.Header.Set("X-ALGOLIA-API-KEY", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		// Add query parameters
		q := httpReq.URL.Query()
		q.Add("indexName", indexName)

		if params.Arguments.Offset != nil {
			q.Add("offset", strconv.Itoa(*params.Arguments.Offset))
		}

		if params.Arguments.Limit != nil {
			q.Add("limit", strconv.Itoa(*params.Arguments.Limit))
		}

		if params.Arguments.Query != nil && *params.Arguments.Query != "" {
			q.Add("query", *params.Arguments.Query)
		}

		httpReq.URL.RawQuery = q.Encode()

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
					Text: fmt.Sprintf("Collections: %v", result),
				},
			},
		}, nil
	})
}
