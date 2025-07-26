package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteIndexParams defines the parameters for deleting an index.
type DeleteIndexParams struct{}

func RegisterDelete(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[DeleteIndexParams]()
	deleteIndexTool := &mcp.Tool{
		Name:        "delete_index",
		Description: "Delete an index by removing all assets and configurations",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteIndexTool, func(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[DeleteIndexParams]) (*mcp.CallToolResultFor[any], error) {
		res, err := index.Delete()
		if err != nil {
			return nil, fmt.Errorf("could not delete index: %v", err)
		}
		return mcputil.JSONToolResult("task", res)
	})
}
