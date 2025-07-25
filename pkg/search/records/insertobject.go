package records

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// InsertObjectParams defines the parameters for inserting an object.
type InsertObjectParams struct {
	Object string `json:"object" jsonschema:"The object to insert or update as a JSON string (must include an objectID field)"`
}

func RegisterInsertObject(mcps *mcp.Server, writeIndex *search.Index) {
	schema, _ := jsonschema.For[InsertObjectParams]()
	insertObjectTool := &mcp.Tool{
		Name:        "insert_object",
		Description: "Insert or update an object in the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, insertObjectTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[InsertObjectParams]) (*mcp.CallToolResultFor[any], error) {
		if writeIndex == nil {
			return nil, fmt.Errorf("write API key not set, cannot insert objects")
		}

		objStr := params.Arguments.Object
		if objStr == "" {
			return nil, fmt.Errorf("invalid object format, expected JSON string")
		}

		// Parse the JSON string into an object
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(objStr), &obj); err != nil {
			return nil, fmt.Errorf("invalid JSON: %v", err)
		}

		// Check if objectID is provided
		if _, exists := obj["objectID"]; !exists {
			return nil, fmt.Errorf("object must include an objectID field")
		}

		// Save the object to the index
		res, err := writeIndex.SaveObject(obj)
		if err != nil {
			return nil, fmt.Errorf("could not save object: %w", err)
		}

		return mcputil.JSONToolResult("insert result", res)
	})
}
