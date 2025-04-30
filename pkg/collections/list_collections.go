package collections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterListCollections registers the list_collections tool with the MCP server.
func RegisterListCollections(mcps *server.MCPServer) {
	listCollectionsTool := mcp.NewTool(
		"collections_list_collections",
		mcp.WithDescription("Retrieve a list of all collections"),
		mcp.WithString(
			"indexName",
			mcp.Description("Name of the index"),
			mcp.Required(),
		),
		mcp.WithNumber(
			"offset",
			mcp.Description("Number of items to skip (default to 0)"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Number of items per fetch (defaults to 10)"),
		),
		mcp.WithString(
			"query",
			mcp.Description("Query to filter collections"),
		),
	)

	mcps.AddTool(listCollectionsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		indexName, _ := req.Params.Arguments["indexName"].(string)
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

		if offset, ok := req.Params.Arguments["offset"].(float64); ok {
			q.Add("offset", strconv.FormatInt(int64(offset), 10))
		}

		if limit, ok := req.Params.Arguments["limit"].(float64); ok {
			q.Add("limit", strconv.FormatInt(int64(limit), 10))
		}

		if query, ok := req.Params.Arguments["query"].(string); ok && query != "" {
			q.Add("query", query)
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

		return mcputil.JSONToolResult("Collections", result)
	})
}
