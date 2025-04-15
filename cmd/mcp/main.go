package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/search/indices"
	"github.com/algolia/mcp/pkg/search/query"
	"github.com/algolia/mcp/pkg/search/records"
	"github.com/algolia/mcp/pkg/search/rules"
	"github.com/algolia/mcp/pkg/search/synonyms"
)

func main() {
	log.Printf("Starting algolia MCP server...")

	var algoliaAppID, algoliaAPIKey, algoliaIndexName string
	if algoliaAppID = os.Getenv("ALGOLIA_APP_ID"); algoliaAppID == "" {
		log.Fatal("ALGOLIA_APP_ID is required")
	}
	if algoliaAPIKey = os.Getenv("ALGOLIA_API_KEY"); algoliaAPIKey == "" {
		log.Fatal("ALGOLIA_API_KEY is required")
	}
	if algoliaIndexName = os.Getenv("ALGOLIA_INDEX_NAME"); algoliaIndexName == "" {
		log.Fatal("ALGOLIA_INDEX_NAME is required")
	}

	client := search.NewClient(algoliaAppID, algoliaAPIKey)
	index := client.InitIndex(algoliaIndexName)

	key, err := client.GetAPIKey(algoliaAPIKey)
	if err != nil {
		log.Fatal("Unable to retrieve ACLs")
	}

	log.Printf("Algolia App ID: %q", algoliaAppID)
	log.Printf("Algolia Index Name: %q", algoliaIndexName)

	log.Printf("Heads up! This MCP server runs with those ACLs: %v.", key.ACL)

	mcps := server.NewMCPServer(
		"algolia-mcp",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// SEARCH TOOLS
	// Tools for managing indices
	indices.RegisterGetSettings(mcps, key.ACL, index)

	// Tools for managing records
	records.RegisterGetObject(mcps, key.ACL, index)
	records.RegisterInsertObject(mcps, key.ACL, index)
	records.RegisterInsertObjects(mcps, key.ACL, index)

	// Tools for searching
	query.RegisterRunQuery(mcps, key.ACL, client, index)

	// Tools for managing rules
	rules.RegisterSearchRules(mcps, key.ACL, index)

	// Tools for managing synonyms
	synonyms.RegisterSearchSynonym(mcps, key.ACL, index)

	if err := server.ServeStdio(mcps); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
