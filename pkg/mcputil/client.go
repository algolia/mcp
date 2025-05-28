package mcputil

import (
	"os"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// JSONToolResult is a convenience method that creates a named JSON-encoded MCP tool result
// from a Go value.
func clientInstance() (*search.Client, *search.Index) {
	var searchClient *search.Client
	var searchIndex *search.Index

	appID := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")
	indexName := os.Getenv("ALGOLIA_INDEX_NAME")

	client = search.NewClient(appID, apiKey)
	index = client.InitIndex(indexName)

	return client, index
}
