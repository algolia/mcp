package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetRecommendStatus registers the get_recommend_status tool with the MCP server.
func RegisterGetRecommendStatus(mcps *server.MCPServer) {
	getRecommendStatusTool := mcp.NewTool(
		"recommend_get_recommend_status",
		mcp.WithDescription("Check the status of a given task"),
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
		mcp.WithNumber(
			"taskID",
			mcp.Description("Unique task identifier"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getRecommendStatusTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		taskIDFloat, ok := req.Params.Arguments["taskID"].(float64)
		if !ok {
			return nil, fmt.Errorf("taskID parameter is required and must be a number")
		}
		taskID := int64(taskIDFloat)

		// Create HTTP client and request
		client := &http.Client{}
		url := fmt.Sprintf("https://%s.algolia.net/1/indexes/%s/%s/task/%s", appID, indexName, model, strconv.FormatInt(taskID, 10))
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

		return mcputil.JSONToolResult("Recommend Task Status", result)
	})
}
