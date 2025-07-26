package records

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetObjectParams defines the parameters for getting an object.
type GetObjectParams struct {
	ObjectID string `json:"objectID" jsonschema:"The object ID to look up"`
}

func RegisterGetObject(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[GetObjectParams]()
	getObjectTool := &mcp.Tool{
		Name:        "get_object",
		Description: "Get an object by its object ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getObjectTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[GetObjectParams]) (*mcp.CallToolResultFor[any], error) {
		objectID := params.Arguments.ObjectID

		var x map[string]any
		if err := index.GetObject(objectID, &x); err != nil {
			return nil, fmt.Errorf("could not get object: %v", err)
		}
		return mcputil.JSONToolResult("object", x)
	})
}
