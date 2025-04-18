package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func RegisterRunQuery(mcps *server.MCPServer, client *search.Client, index *search.Index) {
	runQueryTool := mcp.NewTool(
		"run_query",
		mcp.WithDescription("Run a query against the Algolia search index with advanced options"),
		mcp.WithString(
			"query",
			mcp.Description("The query to run against the index"),
			mcp.Required(),
		),
		mcp.WithString(
			"indexName",
			mcp.Description("The index to search into"),
		),
		mcp.WithNumber(
			"hitsPerPage",
			mcp.Description("The number of hits to return per page"),
		),
		mcp.WithNumber(
			"page",
			mcp.Description("The page number (0-based) to retrieve"),
		),
		mcp.WithString(
			"filters",
			mcp.Description("The filter expression using Algolia's filter syntax (e.g., 'category:Book AND price < 100')"),
		),
		mcp.WithString(
			"facets",
			mcp.Description("Comma-separated list of attributes to facet on"),
		),
		mcp.WithString(
			"restrictSearchableAttributes",
			mcp.Description("Comma-separated list of attributes to search in"),
		),
	)

	mcps.AddTool(runQueryTool, func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName, _ := req.Params.Arguments["indexName"].(string)
		query, _ := req.Params.Arguments["query"].(string)

		opts := []any{}

		// Pagination
		if hitsPerPage, ok := req.Params.Arguments["hitsPerPage"].(float64); ok {
			opts = append(opts, opt.HitsPerPage(int(hitsPerPage)))
		}
		if page, ok := req.Params.Arguments["page"].(float64); ok {
			opts = append(opts, opt.Page(int(page)))
		}

		// Filtering and Faceting
		if filters, ok := req.Params.Arguments["filters"].(string); ok && filters != "" {
			opts = append(opts, opt.Filters(filters))
		}
		if facets, ok := req.Params.Arguments["facets"].(string); ok && facets != "" {
			facetList := strings.Split(facets, ",")
			for i := range facetList {
				facetList[i] = strings.TrimSpace(facetList[i])
			}
			opts = append(opts, opt.Facets(facetList...))
		}

		// Relevance Configuration
		if attrs, ok := req.Params.Arguments["restrictSearchableAttributes"].(string); ok && attrs != "" {
			attrList := strings.Split(attrs, ",")
			for i := range attrList {
				attrList[i] = strings.TrimSpace(attrList[i])
			}
			opts = append(opts, opt.RestrictSearchableAttributes(attrList...))
		}

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

		// Marshal the response directly to JSON
		jsonData, err := json.Marshal(resp["hits"])
		if err != nil {
			return nil, fmt.Errorf("could not marshal search results to JSON: %w", err)
		}

		// Construct the result with the JSON string in the appropriate structure
		// Assuming CallToolResult has Content []mcp.Content
		// AND Assuming mcp.TextContent struct { Type string; Text string } implements mcp.Content
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				// Use mcp.TextContent struct and explicitly set Type
				&mcp.TextContent{
					Type: "text", // Set type explicitly
					Text: string(jsonData),
				},
			},
		}, nil
	})
}
