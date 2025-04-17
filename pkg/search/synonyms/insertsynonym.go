package synonyms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterInsertSynonym(mcps *server.MCPServer, client *search.APIClient, indexName string) {
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
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objectID format"), nil
		}

		synonymStr, ok := req.Params.Arguments["synonym"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid synonym format"), nil
		}

		// Parse synonym
		var synonym search.SynonymHit
		if err := json.Unmarshal([]byte(synonymStr), &synonym); err != nil {
			return nil, fmt.Errorf("could not unmarshal synonym: %w", err)
		}

		result, err := client.SaveSynonym(client.NewApiSaveSynonymRequest(indexName, objectID, &synonym))
		if err != nil {
			return nil, fmt.Errorf("could not save synonym: %w", err)
		}

		return mcputil.JSONToolResult("task", result)
	})
}
