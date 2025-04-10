package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func RegisterInsertObject(mcps *server.MCPServer, writeIndex *search.Index) {
	insertObjectTool := mcp.NewTool(
		"insert_object",
		mcp.WithDescription("Insert or update an object in the Algolia index"),
		mcp.WithString(
			"object",
			mcp.Description("The object to insert or update as a JSON string (must include an objectID field)"),
			mcp.Required(),
		),
	)

	mcps.AddTool(insertObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if writeIndex == nil {
			return mcp.NewToolResultError("write API key not set, cannot insert objects"), nil
		}

		objStr, ok := req.Params.Arguments["object"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid object format, expected JSON string"), nil
		}

		// Parse the JSON string into an object
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(objStr), &obj); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid JSON: %v", err)), nil
		}

		// Check if objectID is provided
		if _, exists := obj["objectID"]; !exists {
			return mcp.NewToolResultError("object must include an objectID field"), nil
		}

		// Save the object to the index
		res, err := writeIndex.SaveObject(obj)
		if err != nil {
			return nil, fmt.Errorf("could not save object: %w", err)
		}

		return jsonResponse("insert result", res)
	})
}

func RegisterInsertObjects(mcps *server.MCPServer, writeIndex *search.Index) {
	insertObjectsTool := mcp.NewTool(
		"insert_objects",
		mcp.WithDescription("Insert or update multiple objects in the Algolia index"),
		mcp.WithString(
			"objects",
			mcp.Description("Array of objects to insert or update as a JSON string (each must include an objectID field)"),
			mcp.Required(),
		),
	)

	mcps.AddTool(insertObjectsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if writeIndex == nil {
			return mcp.NewToolResultError("write API key not set, cannot insert objects"), nil
		}

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

		return jsonResponse("batch insert result", res)
	})
}
