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

func RegisterInsertObject(mcps *server.MCPServer, writeClient *search.Client, writeIndex *search.Index) {
	insertObjectTool := mcp.NewTool(
		"insert_object",
		mcp.WithDescription("Insert or update an object in the Algolia index"),
		mcp.WithString(
			"object",
			mcp.Description("The object to insert or update as a JSON string (must include an objectID field)"),
			mcp.Required(),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("The index to insert objects into"),
		),
	)

	mcps.AddTool(insertObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if writeIndex == nil || writeClient == nil {
			return mcp.NewToolResultError("write API key not set, cannot insert objects"), nil
		}
		writeIndex = mcputil.Index(writeClient, writeIndex, req)

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

		return mcputil.JSONToolResult("insert result", res)
	})
}
