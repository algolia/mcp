package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func RegisterRunQuery(mcps *server.MCPServer, index *search.Index) {
	runQueryTool := mcp.NewTool(
		"run_query",
		mcp.WithDescription("Run a query against the Algolia search index"),
		mcp.WithString(
			"query",
			mcp.Description("The query to run against the index"),
			mcp.Required(),
		),
		mcp.WithNumber(
			"hitsPerPage",
			mcp.Description("The number of hits to return per page"),
		),
	)

	mcps.AddTool(runQueryTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// indexName, _ := req.Params.Arguments["index"].(string)
		query, _ := req.Params.Arguments["query"].(string)

		opts := []any{}
		if hitsPerPage, ok := req.Params.Arguments["hitsPerPage"].(float64); ok {
			opts = append(opts, opt.HitsPerPage(int(hitsPerPage)))
		}

		start := time.Now()
		resp, err := index.Search(query, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not search: %w", err)
		}
		log.Printf("Search for %q took %v", query, time.Since(start))

		return jsonResponse("query results", resp)
	})
}
