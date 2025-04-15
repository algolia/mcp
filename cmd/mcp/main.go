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

	var algoliaAppID, algoliaAPIKey, algoliaIndexName, algoliaWriteAPIKey string
	if algoliaAppID = os.Getenv("ALGOLIA_APP_ID"); algoliaAppID == "" {
		log.Fatal("ALGOLIA_APP_ID is required")
	}
	if algoliaAPIKey = os.Getenv("ALGOLIA_API_KEY"); algoliaAPIKey == "" {
		log.Fatal("ALGOLIA_API_KEY is required")
	}
	if algoliaIndexName = os.Getenv("ALGOLIA_INDEX_NAME"); algoliaIndexName == "" {
		log.Fatal("ALGOLIA_INDEX_NAME is required")
	}

	algoliaWriteAPIKey = os.Getenv("ALGOLIA_WRITE_API_KEY")

	client := search.NewClient(algoliaAppID, algoliaAPIKey)
	index := client.InitIndex(algoliaIndexName)

	log.Printf("Algolia App ID: %q", algoliaAppID)
	log.Printf("Algolia Index Name: %q", algoliaIndexName)

	var writeClient *search.Client
	var writeIndex *search.Index

	if algoliaWriteAPIKey != "" {
		writeClient = search.NewClient(algoliaAppID, algoliaWriteAPIKey)
		writeIndex = writeClient.InitIndex(algoliaIndexName)
		log.Printf("Heads up! This MCP has write capabilities enabled.")
	}

	mcps := server.NewMCPServer(
		"algolia-mcp",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// SEARCH TOOLS
	// Tools for managing indices
	indices.RegisterClear(mcps, writeIndex)
	indices.RegisterGetSettings(mcps, writeIndex)
	indices.RegisterSetSettings(mcps, writeIndex)

	// Tools for managing records
	records.RegisterDeleteObject(mcps, writeIndex)
	records.RegisterGetObject(mcps, index)
	records.RegisterInsertObject(mcps, writeIndex)
	records.RegisterInsertObjects(mcps, writeIndex)

	// Tools for searching
	query.RegisterRunQuery(mcps, client, index)

	// Tools for managing rules
	rules.RegisterDeleteRule(mcps, writeIndex)
	rules.RegisterSearchRules(mcps, index)

	// Tools for managing synonyms
	synonyms.RegisterDeleteSynonym(mcps, writeIndex)
	synonyms.RegisterSearchSynonym(mcps, index)

	if err := server.ServeStdio(mcps); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
