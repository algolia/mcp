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

// RegisterBatchRecommendRules registers the batch_recommend_rules tool with the MCP server.
func RegisterBatchRecommendRules(mcps *server.MCPServer) {
	batchRecommendRulesTool := mcp.NewTool(
		"recommend_batch_recommend_rules",
		mcp.WithDescription("Create or update a batch of Recommend Rules"),
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
			"rules",
			mcp.Description("JSON array of Recommend rules to create or update"),
			mcp.Required(),
		),
		mcp.WithBoolean(
			"clearExistingRules",
			mcp.Description("Whether to replace all existing rules with the provided batch"),
		),
	)

	mcps.AddTool(batchRecommendRulesTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for creating/updating rules
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
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

		rulesJSON, _ := req.Params.Arguments["rules"].(string)
		if rulesJSON == "" {
			return nil, fmt.Errorf("rules parameter is required")
		}

		// Parse rules JSON
		var rules []any
		if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
			return nil, fmt.Errorf("invalid rules JSON: %w", err)
		}

		// Convert rules to JSON
		jsonBody, err := json.Marshal(rules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal rules: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://%s.algolia.net/1/indexes/%s/%s/recommend/rules/batch", appID, indexName, model)
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		httpReq.Header.Set("x-algolia-application-id", appID)
		httpReq.Header.Set("x-algolia-api-key", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		// Add query parameters
		if clearExistingRules, ok := req.Params.Arguments["clearExistingRules"].(bool); ok && clearExistingRules {
			q := httpReq.URL.Query()
			q.Add("clearExistingRules", "true")
			httpReq.URL.RawQuery = q.Encode()
		}

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

		return mcputil.JSONToolResult("Recommend Rules Batch", result)
	})
}
