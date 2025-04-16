package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetSynonym(mcps *server.MCPServer, index *search.Index) {
	getSynonymTool := mcp.NewTool(
		"get_synonym",
		mcp.WithDescription("Get a synonym from the Algolia index by its ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The unique identifier of the synonym to retrieve"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getSynonymTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objectID format"), nil
		}

		synonym, err := index.GetSynonym(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not get synonym: %w", err)
		}

		return mcputil.JSONToolResult("synonym", synonym)
	})
}
