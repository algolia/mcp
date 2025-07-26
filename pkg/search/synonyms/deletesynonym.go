package synonyms

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteSynonymParams defines the parameters for deleting a synonym.
type DeleteSynonymParams struct {
	ObjectID string `json:"objectID" jsonschema:"The object ID to delete"`
}

func RegisterDeleteSynonym(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[DeleteSynonymParams]()
	deleteSynonymTool := &mcp.Tool{
		Name:        "delete_synonym",
		Description: "Delete a synonym by its object ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteSynonymTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteSynonymParams]) (*mcp.CallToolResultFor[any], error) {
		objectID := params.Arguments.ObjectID
		if objectID == "" {
			return nil, fmt.Errorf("invalid object format, expected JSON string")
		}

		resp, err := index.DeleteSynonym(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not delete synonyms: %w", err)
		}

		return mcputil.JSONToolResult("synonym", resp)
	})
}
