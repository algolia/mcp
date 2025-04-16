package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/mcp/pkg/mcputil"
)

const (
	analyticsBaseURL = "https://analytics.algolia.com/2"
)

func RegisterGetAnalytics(mcps *server.MCPServer, index *search.Index, appID, apiKey string) {
	getAnalyticsTool := mcp.NewTool(
		"get_analytics",
		mcp.WithDescription("Fetch analytics data from Algolia"),
		mcp.WithString(
			"startDate",
			mcp.Description("Start date for analytics in YYYY-MM-DD format"),
			mcp.Required(),
		),
		mcp.WithString(
			"endDate",
			mcp.Description("End date for analytics in YYYY-MM-DD format"),
			mcp.Required(),
		),
		mcp.WithString(
			"type",
			mcp.Description("Type of analytics to fetch (searches, clicks, conversions, or topSearches)"),
			mcp.Required(),
			mcp.Enum("searches", "clicks", "conversions", "topSearches"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Maximum number of results to return (for topSearches)"),
		),
	)

	mcps.AddTool(getAnalyticsTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName := index.GetName()
		startDateStr, ok := req.Params.Arguments["startDate"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid startDate format"), nil
		}

		endDateStr, ok := req.Params.Arguments["endDate"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid endDate format"), nil
		}

		analyticsType, ok := req.Params.Arguments["type"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid type format"), nil
		}

		// Parse dates
		startDate, parseErr := time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			return mcp.NewToolResultError("invalid startDate format, expected YYYY-MM-DD"), nil
		}

		endDate, parseErr := time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			return mcp.NewToolResultError("invalid endDate format, expected YYYY-MM-DD"), nil
		}

		// Validate date range
		if endDate.Before(startDate) {
			return mcp.NewToolResultError("endDate must be after startDate"), nil
		}

		// Create HTTP client
		httpClient := &http.Client{}

		// Build request URL based on analytics type
		var url string
		switch analyticsType {
		case "searches":
			url = fmt.Sprintf("%s/searches?index=%s&startDate=%s&endDate=%s", analyticsBaseURL, indexName, startDateStr, endDateStr)
		case "clicks":
			url = fmt.Sprintf("%s/clicks?index=%s&startDate=%s&endDate=%s", analyticsBaseURL, indexName, startDateStr, endDateStr)
		case "conversions":
			url = fmt.Sprintf("%s/conversions?index=%s&startDate=%s&endDate=%s", analyticsBaseURL, indexName, startDateStr, endDateStr)
		case "topSearches":
			limit := 10 // Default limit
			if limitVal, ok := req.Params.Arguments["limit"].(float64); ok {
				limit = int(limitVal)
			}
			url = fmt.Sprintf("%s/top-searches?index=%s&startDate=%s&endDate=%s&limit=%d", analyticsBaseURL, indexName, startDateStr, endDateStr, limit)
		default:
			return mcp.NewToolResultError("invalid analytics type"), nil
		}

		// Create request
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create request: %w", err)
		}

		// Add authentication headers
		httpReq.Header.Set("X-Algolia-Application-Id", appID)
		httpReq.Header.Set("X-Algolia-API-Key", apiKey)

		// Send request
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("could not send request: %w", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
		}

		// Parse response
		var result interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("could not decode response: %w", err)
		}

		return mcputil.JSONToolResult("analytics", result)
	})
}
