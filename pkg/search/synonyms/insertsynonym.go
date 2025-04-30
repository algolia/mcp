package synonyms

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
	synonymsBaseURL = "https://%s.algolia.net/1/indexes/%s/synonyms/%s"
)

func RegisterInsertSynonym(mcps *server.MCPServer, index *search.Index, appID, apiKey string) {
	insertSynonymTool := mcp.NewTool(
		"save_synonym",
		mcp.WithDescription("Save or update a synonym in the Algolia index"),
		mcp.WithString(
			"objectID",
			mcp.Description("The unique identifier of the synonym"),
			mcp.Required(),
		),
		mcp.WithString(
			"synonym",
			mcp.Description("The synonym object as a JSON string. Example schema: {\"objectID\":\"unique_id\",\"type\":\"synonym\",\"synonyms\":[\"word1\",\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"oneWaySynonym\",\"input\":\"word1\",\"synonyms\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"altCorrection1\",\"word\":\"word1\",\"corrections\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"altCorrection2\",\"word\":\"word1\",\"corrections\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"placeholder\",\"placeholder\":\"<em>`,\"replacements\":[\"word1\",\"word2\"]}"),
			mcp.Required(),
		),
	)

	mcps.AddTool(insertSynonymTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName := index.GetName()
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objectID format"), nil
		}

		synonymStr, ok := req.Params.Arguments["synonym"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid synonym format"), nil
		}

		// Parse synonym
		synonym := struct {
			ObjectID string `json:"objectID"`
			Synonym  string `json:"synonym"`
		}{}
		if err := json.Unmarshal([]byte(synonymStr), &synonym); err != nil {
			return nil, fmt.Errorf("could not unmarshal synonym: %w", err)
		}

		// Create HTTP client
		httpClient := &http.Client{}

		// Build request URL
		url := fmt.Sprintf(synonymsBaseURL, appID, indexName, objectID)

		// Create request
		httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(synonymStr)))
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
