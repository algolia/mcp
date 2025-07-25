package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SetSettingsParams defines the parameters for setting index settings.
type SetSettingsParams struct {
	Object string `json:"object" jsonschema:"The settings object as a JSON string"`
}

func RegisterSetSettings(mcps *mcp.Server, writeIndex *search.Index) {
	schema, _ := jsonschema.For[SetSettingsParams]()
	setSettingTool := &mcp.Tool{
		Name:        "set_settings",
		Description: "Change the settings for the Algolia index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, setSettingTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[SetSettingsParams]) (*mcp.CallToolResultFor[any], error) {
		if writeIndex == nil {
			return nil, fmt.Errorf("write API key not set, cannot insert objects")
		}

		objStr := params.Arguments.Object
		if objStr == "" {
			return nil, fmt.Errorf("invalid object format, expected JSON string")
		}

		// Parse the JSON string into an object
		var settings search.Settings
		if err := settings.UnmarshalJSON([]byte(objStr)); err != nil {
			return nil, fmt.Errorf("could not parse settings: %w", err)
		}

		// Save the settings to the index
		res, err := writeIndex.SetSettings(settings)
		if err != nil {
			return nil, fmt.Errorf("could not save object: %w", err)
		}

		return mcputil.JSONToolResult("insert result", res)
	})
}
