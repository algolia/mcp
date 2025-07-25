package rules

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteRuleParams defines the parameters for deleting a rule.
type DeleteRuleParams struct {
	ObjectID string `json:"objectID" jsonschema:"The object ID to delete"`
}

func RegisterDeleteRule(mcps *mcp.Server, index *search.Index) {
	schema, _ := jsonschema.For[DeleteRuleParams]()
	deleteRuleTool := &mcp.Tool{
		Name:        "delete_rule",
		Description: "Delete a rule by its object ID",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, deleteRuleTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteRuleParams]) (*mcp.CallToolResultFor[any], error) {
		objectID := params.Arguments.ObjectID
		if objectID == "" {
			return nil, fmt.Errorf("invalid object format, expected JSON string")
		}

		resp, err := index.DeleteRule(objectID)
		if err != nil {
			return nil, fmt.Errorf("could not delete rule: %w", err)
		}

		return mcputil.JSONToolResult("rule", resp)
	})
}
