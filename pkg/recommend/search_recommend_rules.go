package recommend

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

// RegisterSearchRecommendRules registers the search_recommend_rules tool with the MCP server.
func RegisterSearchRecommendRules(mcps *server.MCPServer) {
	searchRecommendRulesTool := mcp.NewTool(
		"recommend_search_recommend_rules",
		mcp.WithDescription("Search for Recommend rules. Use an empty query to list all rules for this recommendation scenario."),
		mcp.WithString(
			"indexName",
			mcp.Description("Name of the index on which to perform the operation"),
			mcp.Required(),
		),
		mcp.WithString(
			"model",
			mcp.Description("Recommend model (related-products, bought-together, trending-facets, trending-items)"),
			mcp.Required(),
		),
		mcp.WithString(
			"query",
			mcp.Description("Search query"),
		),
		mcp.WithString(
			"context",
			mcp.Description("Only search for rules with matching context"),
		),
		mcp.WithNumber(
			"page",
			mcp.Description("Requested page of the API response"),
		),
		mcp.WithNumber(
			"hitsPerPage",
			mcp.Description("Maximum number of hits per page"),
		),
		mcp.WithBoolean(
			"enabled",
			mcp.Description("Whether to only show rules where the value of their 'enabled' property matches this parameter"),
		),
		mcp.WithString(
			"filters",
			mcp.Description("Filter expression. This only searches for rules matching the filter expression"),
		),
	)

	mcps.AddTool(searchRecommendRulesTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		model, _ := req.Params.Arguments["model"].(string)
		if model == "" {
			return nil, fmt.Errorf("model parameter is required")
		}

		// Prepare request body
		requestBody := make(map[string]any)

		if query, ok := req.Params.Arguments["query"].(string); ok && query != "" {
			requestBody["query"] = query
		}

		if context, ok := req.Params.Arguments["context"].(string); ok && context != "" {
			requestBody["context"] = context
		}

		if page, ok := req.Params.Arguments["page"].(float64); ok {
			requestBody["page"] = int(page)
		}

		if hitsPerPage, ok := req.Params.Arguments["hitsPerPage"].(float64); ok {
			requestBody["hitsPerPage"] = int(hitsPerPage)
		}

		if enabled, ok := req.Params.Arguments["enabled"].(bool); ok {
			requestBody["enabled"] = enabled
		}

		if filters, ok := req.Params.Arguments["filters"].(string); ok && filters != "" {
			requestBody["filters"] = filters
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://%s.algolia.net/1/indexes/%s/%s/recommend/rules/search", appID, indexName, model)
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
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

		return mcputil.JSONToolResult("Recommend Rules Search", result)
	})
}
