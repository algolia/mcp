package mcputil

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// JSONToolResult is a convenience method that creates a named JSON-encoded MCP tool result
// from a Go value.
func JSONToolResult(name string, x any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(x)
	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %w", err)
	}
	return mcp.NewToolResultResource(
		name,
		mcp.TextResourceContents{
			MIMEType: "application/json",
			Text:     string(b),
		},
	), nil
}

// JSONResource is a convenience method that creates a JSON-encoded MCP resource.
func JSONResource(x any) ([]mcp.ResourceContents, error) {
	b, err := json.Marshal(x)
	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %w", err)
	}
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			MIMEType: "application/json",
			Text:     string(b),
		},
	}, nil
}
