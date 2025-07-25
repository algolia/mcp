package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ClearSynonymsParams defines the parameters for clearing synonyms.
type ClearSynonymsParams struct{}

func RegisterClearSynonyms(mcps *mcp.Server, writeIndex *search.Index) {
	schema, _ := jsonschema.For[ClearSynonymsParams]()
	clearSynonymsTool := &mcp.Tool{
		Name:        "clear_synonyms",
		Description: "Clear all synonyms from the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, clearSynonymsTool, func(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[ClearSynonymsParams]) (*mcp.CallToolResultFor[any], error) {
		if writeIndex == nil {
			return nil, fmt.Errorf("write API key not set, cannot clear synonyms")
		}

		res, err := writeIndex.ClearSynonyms()
		if err != nil {
			return nil, fmt.Errorf("could not clear synonyms: %w", err)
		}

		return mcputil.JSONToolResult("clear result", res)
	})
}
