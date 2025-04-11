package mcputil

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// Index is a convenience method that returns the index specified in the request,
// or defaultIndex if no index is specified.
func Index(client *search.Client, defaultIndex *search.Index, req mcp.CallToolRequest) *search.Index {
	indexName, ok := req.Params.Arguments["indexName"].(string)
	if !ok {
		return defaultIndex
	}
	return client.InitIndex(indexName)
}
