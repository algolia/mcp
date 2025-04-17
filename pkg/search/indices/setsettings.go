package indices

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterSetSettings(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	setSettingTool := mcp.NewTool(
		"set_settings",
		mcp.WithDescription("Change the settings for the Algolia index"),
		mcp.WithString(
			"object",
			mcp.Description("The object to insert or update as a JSON string"),
			mcp.Required(),
		),
	)

	mcps.AddTool(setSettingTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if client == nil {
			return mcp.NewToolResultError("write API key not set, cannot insert objects"), nil
		}

		objStr, ok := req.Params.Arguments["object"].(string)
		if !ok {
			return mcp.NewToolResultError("invalid object format, expected JSON string"), nil
		}

		// Parse the JSON string into an object
		var settings *search.IndexSettings
		if err := json.Unmarshal([]byte(objStr), &settings); err != nil {
			return nil, fmt.Errorf("could not parse settings: %w", err)
		}

		// Save the settings to the index
		res, err := client.SetSettings(client.NewApiSetSettingsRequest(indexName, settings))
		if err != nil {
			return nil, fmt.Errorf("could not save object: %w", err)
		}

		return mcputil.JSONToolResult("insert result", res)
	})
}
