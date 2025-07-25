package querysuggestions

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

// UpdateConfigParams defines the parameters for updating a Query Suggestions configuration.
type UpdateConfigParams struct {
	Region                  string  `json:"region" jsonschema:"Analytics region (us or eu)"`
	IndexName               string  `json:"indexName" jsonschema:"Query Suggestions index name"`
	SourceIndices           string  `json:"sourceIndices" jsonschema:"JSON array of source indices configurations"`
	Languages               *string `json:"languages,omitempty" jsonschema:"JSON array of languages or boolean for deduplicating singular and plural suggestions"`
	Exclude                 *string `json:"exclude,omitempty" jsonschema:"JSON array of words or regular expressions to exclude from the suggestions"`
	EnablePersonalization   *bool   `json:"enablePersonalization,omitempty" jsonschema:"Whether to turn on personalized query suggestions"`
	AllowSpecialCharacters  *bool   `json:"allowSpecialCharacters,omitempty" jsonschema:"Whether to include suggestions with special characters"`
}

// RegisterUpdateConfig registers the update_query_suggestions_config tool with the MCP server.
func RegisterUpdateConfig(mcps *mcp.Server) {
	schema, _ := jsonschema.For[UpdateConfigParams]()
	updateConfigTool := &mcp.Tool{
		Name:        "query_suggestions_update_config",
		Description: "Updates a Query Suggestions configuration",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, updateConfigTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[UpdateConfigParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_WRITE_API_KEY") // Note: Using write API key for updating configurations
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_WRITE_API_KEY environment variables are required")
		}

		// Extract parameters
		region := params.Arguments.Region
		if region == "" {
			return nil, fmt.Errorf("region parameter is required")
		}

		indexName := params.Arguments.IndexName
		if indexName == "" {
			return nil, fmt.Errorf("indexName parameter is required")
		}

		sourceIndicesJSON := params.Arguments.SourceIndices
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
		if params.Arguments.Languages != nil && *params.Arguments.Languages != "" {
			var languages any
			if err := json.Unmarshal([]byte(*params.Arguments.Languages), &languages); err != nil {
				return nil, fmt.Errorf("invalid languages JSON: %w", err)
			}
			requestBody["languages"] = languages
		}

		if params.Arguments.Exclude != nil && *params.Arguments.Exclude != "" {
			var exclude []string
			if err := json.Unmarshal([]byte(*params.Arguments.Exclude), &exclude); err != nil {
				return nil, fmt.Errorf("invalid exclude JSON: %w", err)
			}
			requestBody["exclude"] = exclude
		}

		if params.Arguments.EnablePersonalization != nil {
			requestBody["enablePersonalization"] = *params.Arguments.EnablePersonalization
		}

		if params.Arguments.AllowSpecialCharacters != nil {
			requestBody["allowSpecialCharacters"] = *params.Arguments.AllowSpecialCharacters
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Query Suggestions Configuration Updated: %v", result),
				},
			},
		}, nil
	})
}
