package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ClearIndexParams defines the parameters for clearing an index.
type ClearIndexParams struct{}

func RegisterClear(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[ClearIndexParams]()
	clearIndexTool := &mcp.Tool{
		Name:        "clear_index",
		Description: "Clear an index by removing all records",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, clearIndexTool, func(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[ClearIndexParams]) (*mcp.CallToolResultFor[any], error) {
		res, err := index.ClearObjects()
		if err != nil {
			return nil, fmt.Errorf("could not clear index: %v", err)
		}
		return mcputil.JSONToolResult("object", res)
	})
}
