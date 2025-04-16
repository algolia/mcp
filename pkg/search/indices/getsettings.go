package indices

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterGetSettings(mcps *server.MCPServer, index *search.Index) {
	getSettingsTool := mcp.NewTool(
		"get_settings",
		mcp.WithDescription("Get the settings for the Algolia index"),
	)

	mcps.AddTool(getSettingsTool, func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		settings, err := index.GetSettings()
		if err != nil {
			return nil, err
		}
		return mcputil.JSONToolResult("settings", settings)
	})
}
