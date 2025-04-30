package querysuggestions

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

// RegisterUpdateConfig registers the update_query_suggestions_config tool with the MCP server.
func RegisterUpdateConfig(mcps *server.MCPServer) {
	updateConfigTool := mcp.NewTool(
		"query_suggestions_update_config",
		mcp.WithDescription("Updates a Query Suggestions configuration"),
		mcp.WithString(
			"region",
			mcp.Description("Analytics region (us or eu)"),
			mcp.Required(),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("Query Suggestions index name"),
			mcp.Required(),
		),
		mcp.WithString(
			"sourceIndices",
			mcp.Description("JSON array of source indices configurations"),
			mcp.Required(),
		),
		mcp.WithString(
			"languages",
			mcp.Description("JSON array of languages or boolean for deduplicating singular and plural suggestions"),
		),
		mcp.WithString(
			"exclude",
			mcp.Description("JSON array of words or regular expressions to exclude from the suggestions"),
		),
		mcp.WithBoolean(
			"enablePersonalization",
			mcp.Description("Whether to turn on personalized query suggestions"),
		),
		mcp.WithBoolean(
			"allowSpecialCharacters",
			mcp.Description("Whether to include suggestions with special characters"),
		),
	)

	mcps.AddTool(updateConfigTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for updating configurations
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		region, _ := req.Params.Arguments["region"].(string)
		if region == "" {
			return nil, fmt.Errorf("region parameter is required")
		}

		indexName, _ := req.Params.Arguments["indexName"].(string)
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		sourceIndicesJSON, _ := req.Params.Arguments["sourceIndices"].(string)
		if sourceIndicesJSON == "" {
			return nil, fmt.Errorf("sourceIndices parameter is required")
		}

		// Validate region
		if region != "us" && region != "eu" {
			return nil, fmt.Errorf("region must be 'us' or 'eu'")
		}

		// Parse sourceIndices JSON
		var sourceIndices []any
		if err := json.Unmarshal([]byte(sourceIndicesJSON), &sourceIndices); err != nil {
			return nil, fmt.Errorf("invalid sourceIndices JSON: %w", err)
		}

		// Prepare request body
		requestBody := map[string]any{
			"sourceIndices": sourceIndices,
		}

		// Add optional parameters if provided
		if languagesJSON, ok := req.Params.Arguments["languages"].(string); ok && languagesJSON != "" {
			var languages any
			if err := json.Unmarshal([]byte(languagesJSON), &languages); err != nil {
				return nil, fmt.Errorf("invalid languages JSON: %w", err)
			}
			requestBody["languages"] = languages
		}

		if excludeJSON, ok := req.Params.Arguments["exclude"].(string); ok && excludeJSON != "" {
			var exclude []string
			if err := json.Unmarshal([]byte(excludeJSON), &exclude); err != nil {
				return nil, fmt.Errorf("invalid exclude JSON: %w", err)
			}
			requestBody["exclude"] = exclude
		}

		if enablePersonalization, ok := req.Params.Arguments["enablePersonalization"].(bool); ok {
			requestBody["enablePersonalization"] = enablePersonalization
		}

		if allowSpecialCharacters, ok := req.Params.Arguments["allowSpecialCharacters"].(bool); ok {
			requestBody["allowSpecialCharacters"] = allowSpecialCharacters
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://query-suggestions.%s.algolia.com/1/configs/%s", region, indexName)
		httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
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

		return mcputil.JSONToolResult("Query Suggestions Configuration Updated", result)
	})
}
