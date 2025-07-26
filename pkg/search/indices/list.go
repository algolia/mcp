package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListIndicesParams defines the parameters for listing indices.
type ListIndicesParams struct{}

func RegisterList(mcps *mcp.Server, client *search.Client) {
	schema, _ := jsonschema.For[ListIndicesParams]()
	listIndexTool := &mcp.Tool{
		Name:        "list_indices",
		Description: "List the indices in the application",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, listIndexTool, func(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[ListIndicesParams]) (*mcp.CallToolResultFor[any], error) {
		res, err := client.ListIndices()
		if err != nil {
			return nil, fmt.Errorf("could not list indices: %v", err)
		}
		return mcputil.JSONToolResult("indices", res)
	})
}
