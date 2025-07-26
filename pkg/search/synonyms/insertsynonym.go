package synonyms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	synonymsBaseURL = "https://%s.algolia.net/1/indexes/%s/synonyms/%s"
)

// InsertSynonymParams defines the parameters for inserting a synonym.
type InsertSynonymParams struct {
	ObjectID string `json:"objectID" jsonschema:"The unique identifier of the synonym"`
	Synonym  string `json:"synonym" jsonschema:"The synonym object as a JSON string. Example schema: {\"objectID\":\"unique_id\",\"type\":\"synonym\",\"synonyms\":[\"word1\",\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"oneWaySynonym\",\"input\":\"word1\",\"synonyms\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"altCorrection1\",\"word\":\"word1\",\"corrections\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"altCorrection2\",\"word\":\"word1\",\"corrections\":[\"word2\",\"word3\"]} or {\"objectID\":\"unique_id\",\"type\":\"placeholder\",\"placeholder\":\"<em>\",\"replacements\":[\"word1\",\"word2\"]}"`
}

func RegisterInsertSynonym(mcps *mcp.Server, index *search.Index, appID, apiKey string) {
	schema, _ := jsonschema.For[InsertSynonymParams]()
	insertSynonymTool := &mcp.Tool{
		Name:        "save_synonym",
		Description: "Save or update a synonym in the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, insertSynonymTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[InsertSynonymParams]) (*mcp.CallToolResultFor[any], error) {
		indexName := index.GetName()
		objectID := params.Arguments.ObjectID
		if objectID == "" {
			return nil, fmt.Errorf("invalid objectID format")
		}

		synonymStr := params.Arguments.Synonym
		if synonymStr == "" {
			return nil, fmt.Errorf("invalid synonym format")
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
