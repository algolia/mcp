package mcputil

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

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
