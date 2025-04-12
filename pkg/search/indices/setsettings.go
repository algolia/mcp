package indices

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterSetSettings(mcps *server.MCPServer, writeIndex *search.Index) {
	setSettingTool := mcp.NewTool(
		"set_settings",
		mcp.WithDescription("Change the settings for the Algolia index"),
		mcp.WithString(
			"object",
			mcp.Description("The object to insert or update as a JSON string"),
			mcp.Required(),
		),
	)

	mcps.AddTool(setSettingTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Convert the object to search.Settings
		settings := search.Settings{}
		settingsBytes, err := json.Marshal(obj)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal settings: %v", err)), nil
		}
		if err := json.Unmarshal(settingsBytes, &settings); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to unmarshal settings: %v", err)), nil
		}

		// Save the settings to the index
		res, err := writeIndex.SetSettings(settings)
		if err != nil {
			return nil, fmt.Errorf("could not save object: %w", err)
		}

		return mcputil.JSONToolResult("insert result", res)
	})
}
