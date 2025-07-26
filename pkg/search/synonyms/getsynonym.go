package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetSynonymParams defines the parameters for getting a synonym.
type GetSynonymParams struct {
	ObjectID string `json:"objectID" jsonschema:"The unique identifier of the synonym to retrieve"`
}

func RegisterGetSynonym(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[GetSynonymParams]()
	getSynonymTool := &mcp.Tool{
		Name:        "get_synonym",
		Description: "Get a synonym from the Algolia index by its ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getSynonymTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[GetSynonymParams]) (*mcp.CallToolResultFor[any], error) {
		objectID := params.Arguments.ObjectID
		if objectID == "" {
			return nil, fmt.Errorf("invalid objectID format")
		}

		synonym, err := index.GetSynonym(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not get synonym: %w", err)
		}

		return mcputil.JSONToolResult("synonym", synonym)
	})
}
