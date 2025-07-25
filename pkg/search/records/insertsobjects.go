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

// InsertObjectsParams defines the parameters for inserting multiple objects.
type InsertObjectsParams struct {
	Objects string `json:"objects" jsonschema:"Array of objects to insert or update as a JSON string (each must include an objectID field)"`
}

func RegisterInsertObjects(mcps *mcp.Server, writeIndex *search.Index) {
	schema, _ := jsonschema.For[InsertObjectsParams]()
	insertObjectsTool := &mcp.Tool{
		Name:        "insert_objects",
		Description: "Insert or update multiple objects in the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, insertObjectsTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[InsertObjectsParams]) (*mcp.CallToolResultFor[any], error) {
		if writeIndex == nil {
			return nil, fmt.Errorf("write API key not set, cannot insert objects")
		}

		objsStr := params.Arguments.Objects
		if objsStr == "" {
			return nil, fmt.Errorf("invalid objects format, expected JSON string")
		}

		// Parse the JSON string into an array of objects
		var objects []map[string]interface{}
		if err := json.Unmarshal([]byte(objsStr), &objects); err != nil {
			return nil, fmt.Errorf("invalid JSON: %v", err)
		}

		// Check if all objects have an objectID
		for i, obj := range objects {
			if _, exists := obj["objectID"]; !exists {
				return nil, fmt.Errorf("object at index %d must include an objectID field", i)
			}
		}

		// Save the objects to the index
		res, err := writeIndex.SaveObjects(objects)
		if err != nil {
			return nil, fmt.Errorf("could not save objects: %w", err)
		}

		return mcputil.JSONToolResult("batch insert result", res)
	})
}
