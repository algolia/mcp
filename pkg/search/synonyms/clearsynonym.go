package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterClearSynonyms(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	clearSynonymsTool := mcp.NewTool(
		"clear_synonyms",
		mcp.WithDescription("Clear all synonyms from the Algolia index"),
	)

	mcps.AddTool(clearSynonymsTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if client == nil {
			return mcp.NewToolResultError("write API key not set, cannot clear synonyms"), nil
		}

		res, err := client.ClearSynonyms(client.NewApiClearSynonymsRequest(indexName))
		if err != nil {
			return nil, fmt.Errorf("could not clear synonyms: %w", err)
		}

		return mcputil.JSONToolResult("clear result", res)
	})
}
