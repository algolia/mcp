package synonyms

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterDeleteSynonym(mcps *server.MCPServer, index *search.Index) {
	DeleteSynonymTool := mcp.NewTool(
		"delete_synonym",
		mcp.WithDescription("Delete a synonym by its object ID"),
		mcp.WithString(
			"objectID",
			mcp.Description("The object ID to delete"),
			mcp.Required(),
		),
	)

	mcps.AddTool(DeleteSynonymTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectID, ok := req.Params.Arguments["objectID"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid object format, expected JSON string"), nil
		}

		resp, err := index.DeleteRule(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not delete synonyms: %w", err)
		}

		return mcputil.JSONToolResult("synonym", resp)
	})
}
