package query

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterRunQuery(mcps *server.MCPServer, client *search.APIClient, indexName string) {
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
		indexNameOpt, _ := req.Params.Arguments["indexName"].(string)
		query, _ := req.Params.Arguments["query"].(string)

		searchParams := search.NewEmptySearchForHits()
		searchParams.SetIndexName(indexName)
		if indexNameOpt != "" {
			searchParams.SetIndexName(indexNameOpt)
		}

		// Pagination
		if hitsPerPage, ok := req.Params.Arguments["hitsPerPage"].(int32); ok {
			searchParams.SetHitsPerPage(hitsPerPage)
		}
		if page, ok := req.Params.Arguments["page"].(int32); ok {
			searchParams.SetPage(page)
		}

		// Filtering and Faceting
		if filters, ok := req.Params.Arguments["filters"].(string); ok && filters != "" {
			searchParams.SetFilters(filters)
		}
		if facets, ok := req.Params.Arguments["facets"].(string); ok && facets != "" {
			facetList := strings.Split(facets, ",")
			for i := range facetList {
				facetList[i] = strings.TrimSpace(facetList[i])
			}
			searchParams.SetFacets(facetList)
		}

		// Relevance Configuration
		if attrs, ok := req.Params.Arguments["restrictSearchableAttributes"].(string); ok && attrs != "" {
			attrList := strings.Split(attrs, ",")
			searchParams.SetRestrictSearchableAttributes(attrList)
		}

		start := time.Now()
		resp, err := client.Search(client.NewApiSearchRequest(
			search.NewEmptySearchMethodParams().SetRequests(
				[]search.SearchQuery{*search.SearchForHitsAsSearchQuery(searchParams)},
			),
		),
		)
		if err != nil {
			return nil, fmt.Errorf("could not search: %w", err)
		}
		log.Printf("Search for %q took %v", query, time.Since(start))

		return mcputil.JSONToolResult("query results", resp)
	})
}
