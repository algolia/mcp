package mcputil

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// JSONToolResult is a convenience method that creates a named JSON-encoded MCP tool result
// from a Go value.
func JSONToolResult(name string, x any) (*mcp.CallToolResultFor[any], error) {
	b, err := json.Marshal(x)
	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %w", err)
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%s: %s", name, string(b)),
			},
		},
	}, nil
}

// JSONResource is a convenience method that creates a JSON-encoded MCP resource.
func JSONResource(x any) ([]mcp.Content, error) {
	b, err := json.Marshal(x)
	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %w", err)
	}
	return []mcp.Content{
		&mcp.TextContent{
			Text: string(b),
		},
	}, nil
}
