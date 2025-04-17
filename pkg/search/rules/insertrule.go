package rules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

const (
	rulesBaseURL = "https://%s.algolia.net/1/indexes/%s/rules/%s"
)

func RegisterInsertRule(mcps *server.MCPServer, index *search.Index, appID, apiKey string) {
	insertRuleTool := mcp.NewTool(
		"save_rule",
		mcp.WithDescription("Save or update a rule in the Algolia index"),
		mcp.WithString(
			"objectID",
			mcp.Description("The unique identifier of the rule"),
			mcp.Required(),
		),
		mcp.WithString(
			"rule",
			mcp.Description("The rule object as a JSON string. Example schema: {\"objectID\":\"pattern\",\"conditions\":[{\"anchoring\":\"contains\",\"pattern\":\"pattern\"}],\"consequence\":{\"params\":{\"filters\":\"attribute:value\",\"query\":\"\"},\"promote\":[{\"objectID\":\"id\",\"position\":0}]},\"description\":\"description\"}"),
			mcp.Required(),
		),
	)

	mcps.AddTool(insertRuleTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName := index.GetName()
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objectID format"), nil
		}

		ruleStr, ok := req.Params.Arguments["rule"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid rule format, expected JSON string"), nil
		}

		// Create HTTP client
		httpClient := &http.Client{}

		// Build request URL
		url := fmt.Sprintf(rulesBaseURL, appID, indexName, objectID)

		// Create request
		httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(ruleStr))
		if err != nil {
			return nil, fmt.Errorf("could not create request: %w", err)
		}

		// Add headers
		httpReq.Header.Set("Content-Type", "application/json")
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

		return mcputil.JSONToolResult("task", result)
	})
}
