package synonyms

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterSearchSynonym(mcps *server.MCPServer, index *search.Index) {
	searchSynonymTool := mcp.NewTool(
		"search_synonyms",
		mcp.WithDescription("Search for synonyms in the Algolia index that match a query"),
		mcp.WithString(
			"query",
			mcp.Description("The query to find synonyms for"),
		),
	)

	mcps.AddTool(searchSynonymTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.Params.Arguments["query"].(string)

		resp, err := index.SearchSynonyms(query)
		if err != nil {
			return nil, fmt.Errorf("could not search synonyms: %w", err)
		}

		return mcputil.JSONToolResult("synonyms", resp)
	})
}
