package recommend

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

// BatchRecommendRulesParams defines the parameters for batch operations on recommend rules.
type BatchRecommendRulesParams struct {
	IndexName          string `json:"indexName" jsonschema:"Name of the index on which to perform the operation"`
	Model              string `json:"model" jsonschema:"Recommend model (related-products, bought-together, trending-facets, trending-items)"`
	Rules              string `json:"rules" jsonschema:"JSON array of Recommend rules to create or update"`
	ClearExistingRules *bool  `json:"clearExistingRules,omitempty" jsonschema:"Whether to replace all existing rules with the provided batch"`
}

// RegisterBatchRecommendRules registers the batch_recommend_rules tool with the MCP server.
func RegisterBatchRecommendRules(mcps *mcp.Server) {
	schema, _ := jsonschema.For[BatchRecommendRulesParams]()
	batchRecommendRulesTool := &mcp.Tool{
		Name:        "recommend_batch_recommend_rules",
		Description: "Create or update a batch of Recommend Rules",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, batchRecommendRulesTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[BatchRecommendRulesParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for creating/updating rules
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		indexName := params.Arguments.IndexName
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		model := params.Arguments.Model
		if model == "" {
			return nil, fmt.Errorf("model parameter is required")
		}

		rulesJSON := params.Arguments.Rules
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
		if params.Arguments.ClearExistingRules != nil && *params.Arguments.ClearExistingRules {
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Recommend Rules Batch: %v", result),
				},
			},
		}, nil
	})
}
