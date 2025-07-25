package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MoveIndexParams defines the parameters for moving an index.
type MoveIndexParams struct {
	IndexName string `json:"indexName" jsonschema:"The name of the destination index"`
}

func RegisterMove(mcps *mcp.Server, client *search.Client, index *search.Index) {
	schema, _ := jsonschema.For[MoveIndexParams]()
	moveIndexTool := &mcp.Tool{
		Name:        "move_index",
		Description: "Move an index to another index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, moveIndexTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[MoveIndexParams]) (*mcp.CallToolResultFor[any], error) {
		dst := params.Arguments.IndexName
		if dst == "" {
			return nil, fmt.Errorf("invalid indexName format, expected non-empty string")
		}

		res, err := client.CopyIndex(index.GetName(), dst)
		if err != nil {
			return nil, fmt.Errorf("could not move index: %v", err)
		}
		return mcputil.JSONToolResult("task", res)
	})
}
