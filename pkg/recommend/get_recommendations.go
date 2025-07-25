package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/recommend"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetRecommendationsParams defines the parameters for retrieving recommendations.
type GetRecommendationsParams struct {
	Requests string `json:"requests" jsonschema:"JSON array of recommendation requests. Each request must include 'indexName', 'threshold', and a model-specific configuration."`
}

// RegisterGetRecommendations registers the get_recommendations tool with the MCP server.
func RegisterGetRecommendations(mcps *mcp.Server) {
	schema, _ := jsonschema.For[GetRecommendationsParams]()
	getRecommendationsTool := &mcp.Tool{
		Name:        "recommend_get_recommendations",
		Description: "Retrieve recommendations from selected AI models",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getRecommendationsTool, func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetRecommendationsParams]) (*mcp.CallToolResultFor[any], error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		requestsJSON := params.Arguments.Requests
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

		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Recommendations: %v", res),
				},
			},
		}, nil
	})
}
