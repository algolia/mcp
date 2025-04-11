package records

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterInsertObjects(mcps *server.MCPServer, writeClient *search.Client, writeIndex *search.Index) {
	insertObjectsTool := mcp.NewTool(
		"insert_objects",
		mcp.WithDescription("Insert or update multiple objects in the Algolia index"),
		mcp.WithString(
			"objects",
			mcp.Description("Array of objects to insert or update as a JSON string (each must include an objectID field)"),
			mcp.Required(),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("The index to insert objects into"),
		),
	)

	mcps.AddTool(insertObjectsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if writeIndex == nil || writeClient == nil {
			return mcp.NewToolResultError("write API key not set, cannot insert objects"), nil
		}
		writeIndex = mcputil.Index(writeClient, writeIndex, req)

		objsStr, ok := req.Params.Arguments["objects"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid objects format, expected JSON string"), nil
		}

		// Parse the JSON string into an array of objects
		var objects []map[string]interface{}
		if err := json.Unmarshal([]byte(objsStr), &objects); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid JSON: %v", err)), nil
		}

		// Check if all objects have an objectID
		for i, obj := range objects {
			if _, exists := obj["objectID"]; !exists {
				return mcp.NewToolResultError(fmt.Sprintf("object at index %d must include an objectID field", i)), nil
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
