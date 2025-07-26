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

// SearchRecommendRulesParams defines the parameters for searching recommend rules.
type SearchRecommendRulesParams struct {
	IndexName   string  `json:"indexName" jsonschema:"Name of the index on which to perform the operation"`
	Model       string  `json:"model" jsonschema:"Recommend model (related-products, bought-together, trending-facets, trending-items)"`
	Query       *string `json:"query,omitempty" jsonschema:"Search query"`
	Context     *string `json:"context,omitempty" jsonschema:"Only search for rules with matching context"`
	Page        *int    `json:"page,omitempty" jsonschema:"Requested page of the API response"`
	HitsPerPage *int    `json:"hitsPerPage,omitempty" jsonschema:"Maximum number of hits per page"`
	Enabled     *bool   `json:"enabled,omitempty" jsonschema:"Whether to only show rules where the value of their 'enabled' property matches this parameter"`
	Filters     *string `json:"filters,omitempty" jsonschema:"Filter expression. This only searches for rules matching the filter expression"`
}

// RegisterSearchRecommendRules registers the search_recommend_rules tool with the MCP server.
func RegisterSearchRecommendRules(mcps *mcp.Server) {
	schema, _ := jsonschema.For[SearchRecommendRulesParams]()
	searchRecommendRulesTool := &mcp.Tool{
		Name:        "recommend_search_recommend_rules",
		Description: "Search for Recommend rules. Use an empty query to list all rules for this recommendation scenario.",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, searchRecommendRulesTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[SearchRecommendRulesParams]) (*mcp.CallToolResultFor[any], error) {
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

		model := params.Arguments.Model
		if model == "" {
			return nil, fmt.Errorf("model parameter is required")
		}

		// Prepare request body
		requestBody := make(map[string]any)

		if params.Arguments.Query != nil && *params.Arguments.Query != "" {
			requestBody["query"] = *params.Arguments.Query
		}

		if params.Arguments.Context != nil && *params.Arguments.Context != "" {
			requestBody["context"] = *params.Arguments.Context
		}

		if params.Arguments.Page != nil {
			requestBody["page"] = *params.Arguments.Page
		}

		if params.Arguments.HitsPerPage != nil {
			requestBody["hitsPerPage"] = *params.Arguments.HitsPerPage
		}

		if params.Arguments.Enabled != nil {
			requestBody["enabled"] = *params.Arguments.Enabled
		}

		if params.Arguments.Filters != nil && *params.Arguments.Filters != "" {
			requestBody["filters"] = *params.Arguments.Filters
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Recommend Rules Search: %v", result),
				},
			},
		}, nil
	})
}
