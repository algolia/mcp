package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/recommend"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetRecommendations registers the get_recommendations tool with the MCP server.
func RegisterGetRecommendations(mcps *server.MCPServer) {
	getRecommendationsTool := mcp.NewTool(
		"recommend_get_recommendations",
		mcp.WithDescription("Retrieve recommendations from selected AI models"),
		mcp.WithString(
			"requests",
			mcp.Description("JSON array of recommendation requests. Each request must include 'indexName', 'threshold', and a model-specific configuration."),
			mcp.Required(),
		),
	)

	mcps.AddTool(getRecommendationsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		requestsJSON, _ := req.Params.Arguments["requests"].(string)
		if requestsJSON == "" {
			return nil, fmt.Errorf("requests parameter is required")
		}

		// Parse requests JSON
		var rawRequests []map[string]interface{}
		if err := json.Unmarshal([]byte(requestsJSON), &rawRequests); err != nil {
			return nil, fmt.Errorf("invalid requests JSON: %w", err)
		}

		// Convert raw requests to RecommendationsOptions
		var options []recommend.RecommendationsOptions
		for _, rawReq := range rawRequests {
			// Extract required fields
			indexName, _ := rawReq["indexName"].(string)
			if indexName == "" {
				return nil, fmt.Errorf("indexName is required for each request")
			}

			modelStr, _ := rawReq["model"].(string)
			if modelStr == "" {
				return nil, fmt.Errorf("model is required for each request")
			}
			model := recommend.RecommendationModel(modelStr)

			objectID, _ := rawReq["objectID"].(string)
			if objectID == "" {
				return nil, fmt.Errorf("objectID is required for each request")
			}

			thresholdFloat, _ := rawReq["threshold"].(float64)
			if thresholdFloat == 0 {
				return nil, fmt.Errorf("threshold is required for each request")
			}
			threshold := int(thresholdFloat)

			// Create options
			opt := recommend.RecommendationsOptions{
				IndexName: indexName,
				Model:     model,
				ObjectID:  objectID,
				Threshold: threshold,
			}

			// Add optional fields
			if maxRecsFloat, ok := rawReq["maxRecommendations"].(float64); ok {
				maxRecs := int(maxRecsFloat)
				opt.MaxRecommendations = &maxRecs
			}

			// Add query parameters if provided
			if _, ok := rawReq["queryParameters"].(map[string]interface{}); ok {
				// For now, we'll just create an empty QueryParams
				// In a real implementation, you would need to convert the map to the appropriate types
				opt.QueryParameters = &search.QueryParams{}
			}

			options = append(options, opt)
		}

		// Create Algolia Recommend client
		client := recommend.NewClient(appID, apiKey)

		// Get recommendations
		res, err := client.GetRecommendations(options)
		if err != nil {
			return nil, fmt.Errorf("failed to get recommendations: %w", err)
		}

		return mcputil.JSONToolResult("Recommendations", res)
	})
}
