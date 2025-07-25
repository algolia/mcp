package indices

import (
	"context"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetSettingsParams defines the parameters for getting settings.
type GetSettingsParams struct{}

func RegisterGetSettings(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[GetSettingsParams]()
	getSettingsTool := &mcp.Tool{
		Name:        "get_settings",
		Description: "Get the settings for the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, getSettingsTool, func(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[GetSettingsParams]) (*mcp.CallToolResultFor[any], error) {
		settings, err := index.GetSettings()
		if err != nil {
			return nil, err
		}
		return mcputil.JSONToolResult("settings", settings)
	})
}
