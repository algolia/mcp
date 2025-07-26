package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearchSynonymParams defines the parameters for searching synonyms.
type SearchSynonymParams struct {
	Query string `json:"query" jsonschema:"The query to find synonyms for"`
}

func RegisterSearchSynonym(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[SearchSynonymParams]()
	searchSynonymTool := &mcp.Tool{
		Name:        "search_synonyms",
		Description: "Search for synonyms in the Algolia index that match a query",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, searchSynonymTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[SearchSynonymParams]) (*mcp.CallToolResultFor[any], error) {
		query := params.Arguments.Query

		resp, err := index.SearchSynonyms(query)
		if err != nil {
			return nil, fmt.Errorf("could not search synonyms: %w", err)
		}

		return mcputil.JSONToolResult("synonyms", resp)
	})
}
