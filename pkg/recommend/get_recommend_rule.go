package recommend

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

// RegisterGetRecommendRule registers the get_recommend_rule tool with the MCP server.
func RegisterGetRecommendRule(mcps *server.MCPServer) {
	getRecommendRuleTool := mcp.NewTool(
		"recommend_get_recommend_rule",
		mcp.WithDescription("Retrieve a Recommend rule that you previously created in the Algolia dashboard"),
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
			"objectID",
			mcp.Description("Unique record identifier"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getRecommendRuleTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		objectID, _ := req.Params.Arguments["objectID"].(string)
		if objectID == "" {
			return nil, fmt.Errorf("objectID parameter is required")
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://%s.algolia.net/1/indexes/%s/%s/recommend/rules/%s", appID, indexName, model, objectID)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
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

		return mcputil.JSONToolResult("Recommend Rule", result)
	})
}
