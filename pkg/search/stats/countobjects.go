package stats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterCountObjects(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	runQueryTool := mcp.NewTool(
		"count_objects",
		mcp.WithDescription("Count the number of objects in the Algolia search index."),
		mcp.WithString(
			"indexName",
			mcp.Description("The index to count the objects"),
		),
	)

	mcps.AddTool(runQueryTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName, _ := req.Params.Arguments["indexName"].(string)
		query, _ := req.Params.Arguments["query"].(string)

		opts := []any{opt.HitsPerPage(0)}

		currentIndex := index
		if indexName != "" {
			currentIndex = client.InitIndex(indexName)
		}

		start := time.Now()
		resp, err := currentIndex.Search(query, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not search: %w", err)
		}
		log.Printf("Search for %q took %v", query, time.Since(start))

		return mcputil.JSONToolResult("query results", resp.NbHits)
	})
}
