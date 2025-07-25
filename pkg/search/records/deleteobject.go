package records

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteObjectParams defines the parameters for deleting an object.
type DeleteObjectParams struct {
	ObjectID string `json:"objectID" jsonschema:"The object ID to delete"`
}

func RegisterDeleteObject(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[DeleteObjectParams]()
	deleteObjectTool := &mcp.Tool{
		Name:        "delete_object",
		Description: "Delete an object by its object ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteObjectTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteObjectParams]) (*mcp.CallToolResultFor[any], error) {
		objectID := params.Arguments.ObjectID

		res, err := index.DeleteObject(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not delete object: %v", err)
		}
		return mcputil.JSONToolResult("object", res)
	})
}
