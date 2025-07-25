package indices

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CopyIndexParams defines the parameters for copying an index.
type CopyIndexParams struct {
	IndexName string `json:"indexName" jsonschema:"The name of the destination index"`
}

func RegisterCopy(mcps *mcp.Server, client *search.Client, index *search.Index) {
	schema, _ := jsonschema.For[CopyIndexParams]()
	copyIndexTool := &mcp.Tool{
		Name:        "copy_index",
		Description: "Copy an index to a another index",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, copyIndexTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[CopyIndexParams]) (*mcp.CallToolResultFor[any], error) {
		dst := params.Arguments.IndexName
		if dst == "" {
			return nil, fmt.Errorf("invalid indexName format, expected non-empty string")
		}

		res, err := client.CopyIndex(index.GetName(), dst)
		if err != nil {
			return nil, fmt.Errorf("could not copy index: %v", err)
		}
		return mcputil.JSONToolResult("task", res)
	})
}
