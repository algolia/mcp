package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetRecommendRuleParams defines the parameters for retrieving a recommend rule.
type GetRecommendRuleParams struct {
	IndexName string `json:"indexName" jsonschema:"Name of the index on which to perform the operation"`
	Model     string `json:"model" jsonschema:"Recommend model (related-products, bought-together, trending-facets, trending-items)"`
	ObjectID  string `json:"objectID" jsonschema:"Unique record identifier"`
}

// RegisterGetRecommendRule registers the get_recommend_rule tool with the MCP server.
func RegisterGetRecommendRule(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetRecommendRuleParams]()
	getRecommendRuleTool := &mcp.Tool{
		Name:        "recommend_get_recommend_rule",
		Description: "Retrieve a Recommend rule that you previously created in the Algolia dashboard",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getRecommendRuleTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetRecommendRuleParams]) (*mcp.CallToolResultFor[any], error) {
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

		objectID := params.Arguments.ObjectID
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Recommend Rule: %v", result),
				},
			},
		}, nil
	})
}
