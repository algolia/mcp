package indices

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterGetSettings(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	getSettingsTool := mcp.NewTool(
		"get_settings",
		mcp.WithDescription("Get the settings for the Algolia index"),
		mcp.WithString(
			"indexName",
			mcp.Description("The index to retrieve settings from"),
		),
	)

	mcps.AddTool(getSettingsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		index = mcputil.Index(client, index, req)

		settings, err := index.GetSettings()
		if err != nil {
			return nil, err
		}
		return mcputil.JSONToolResult("settings", settings)
	})
}
