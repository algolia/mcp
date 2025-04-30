package abtesting

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

// RegisterEstimateABTest registers the estimate_abtest tool with the MCP server.
func RegisterEstimateABTest(mcps *server.MCPServer) {
	estimateABTestTool := mcp.NewTool(
		"abtesting_estimate_abtest",
		mcp.WithDescription("Estimate the sample size and duration of an A/B test based on historical traffic"),
		mcp.WithString(
			"variants",
			mcp.Description("A/B test variants as JSON array (exactly 2 variants required). Each variant must have 'index' and 'trafficPercentage' fields, and may optionally have 'description' and 'customSearchParameters' fields."),
			mcp.Required(),
		),
		mcp.WithString(
			"configuration",
			mcp.Description("A/B test configuration as JSON object. Must include 'minimumDetectableEffect' with 'size' and 'metric' fields. May optionally include 'outliers' and 'emptySearch' settings."),
			mcp.Required(),
		),
	)

	mcps.AddTool(estimateABTestTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appID := os.Getenv("ALGOLIA_APP_ID")
		apiKey := os.Getenv("ALGOLIA_API_KEY")
		if appID == "" || apiKey == "" {
			return nil, fmt.Errorf("ALGOLIA_APP_ID and ALGOLIA_API_KEY environment variables are required")
		}

		// Extract parameters
		variantsJSON, _ := req.Params.Arguments["variants"].(string)
		configJSON, _ := req.Params.Arguments["configuration"].(string)

		// Parse variants JSON
		var variants []any
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err != nil {
			return nil, fmt.Errorf("invalid variants JSON: %w", err)
		}

		if len(variants) != 2 {
			return nil, fmt.Errorf("exactly 2 variants are required")
		}

		// Parse configuration JSON
		var config map[string]any
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return nil, fmt.Errorf("invalid configuration JSON: %w", err)
		}

		// Check for required minimumDetectableEffect
		if _, ok := config["minimumDetectableEffect"]; !ok {
			return nil, fmt.Errorf("configuration must include 'minimumDetectableEffect'")
		}

		// Prepare request body
		requestBody := map[string]any{
			"configuration": config,
			"variants":      variants,
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create HTTP client and request
		client := &http.Client{}
		url := "https://analytics.algolia.com/2/abtests/estimate"
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

		return mcputil.JSONToolResult("AB Test Estimate", result)
	})
}
