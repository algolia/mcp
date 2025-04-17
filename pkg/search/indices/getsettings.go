package indices

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterGetSettings(mcps *server.MCPServer, client *search.APIClient, indexName string) {
	getSettingsTool := mcp.NewTool(
		"get_settings",
		mcp.WithDescription("Get the settings for the Algolia index"),
	)

	mcps.AddTool(getSettingsTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		settings, err := client.GetSettings(client.NewApiGetSettingsRequest(indexName))
		if err != nil {
			return nil, err
		}
		return mcputil.JSONToolResult("settings", settings)
	})
}
