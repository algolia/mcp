package query

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RunQueryParams defines the parameters for running a query.
type RunQueryParams struct {
	Query                        string `json:"query" jsonschema:"The query to run against the index"`
	IndexName                    string `json:"indexName,omitempty" jsonschema:"The index to search into"`
	HitsPerPage                  *int   `json:"hitsPerPage,omitempty" jsonschema:"The number of hits to return per page"`
	Page                         *int   `json:"page,omitempty" jsonschema:"The page number (0-based) to retrieve"`
	Filters                      string `json:"filters,omitempty" jsonschema:"The filter expression using Algolia's filter syntax (e.g., 'category:Book AND price < 100')"`
	Facets                       string `json:"facets,omitempty" jsonschema:"Comma-separated list of attributes to facet on"`
	RestrictSearchableAttributes string `json:"restrictSearchableAttributes,omitempty" jsonschema:"Comma-separated list of attributes to search in"`
}

func RegisterRunQuery(mcps *mcp.Server, client *search.Client, index *search.Index) {
	schema, _ := jsonschema.For[RunQueryParams]()
	runQueryTool := &mcp.Tool{
		Name:        "run_query",
		Description: "Run a query against the Algolia search index with advanced options",
		InputSchema: schema,
	}

	mcp.AddTool(mcps, runQueryTool, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[RunQueryParams]) (*mcp.CallToolResultFor[any], error) {
		indexName := params.Arguments.IndexName
		query := params.Arguments.Query

		opts := []any{}

		// Pagination
		if params.Arguments.HitsPerPage != nil {
			opts = append(opts, opt.HitsPerPage(*params.Arguments.HitsPerPage))
		}
		if params.Arguments.Page != nil {
			opts = append(opts, opt.Page(*params.Arguments.Page))
		}

		// Filtering and Faceting
		if params.Arguments.Filters != "" {
			opts = append(opts, opt.Filters(params.Arguments.Filters))
		}
		if params.Arguments.Facets != "" {
			facetList := strings.Split(params.Arguments.Facets, ",")
			for i := range facetList {
				facetList[i] = strings.TrimSpace(facetList[i])
			}
			opts = append(opts, opt.Facets(facetList...))
		}

		// Relevance Configuration
		if params.Arguments.RestrictSearchableAttributes != "" {
			attrList := strings.Split(params.Arguments.RestrictSearchableAttributes, ",")
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

		return mcputil.JSONToolResult("query results", resp)
	})
}
