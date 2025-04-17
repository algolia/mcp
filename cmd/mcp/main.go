package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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

	client, err := search.NewClient(algoliaAppID, algoliaAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Algolia client: %v", err)
	}

	log.Printf("Algolia App ID: %q", algoliaAppID)
	log.Printf("Algolia Index Name: %q", algoliaIndexName)

	var writeClient *search.APIClient

	if algoliaWriteAPIKey != "" {
		writeClient, err = search.NewClient(algoliaAppID, algoliaWriteAPIKey)
		if err != nil {
			log.Fatalf("Failed to create Algolia write client: %v", err)
		}
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
	indices.RegisterClear(mcps, writeClient, algoliaIndexName)
	indices.RegisterCopy(mcps, writeClient, algoliaIndexName)
	indices.RegisterGetSettings(mcps, writeClient, algoliaIndexName)
	indices.RegisterList(mcps, writeClient)
	indices.RegisterMove(mcps, writeClient, algoliaIndexName)
	indices.RegisterSetSettings(mcps, writeClient, algoliaIndexName)

	// Tools for managing records
	records.RegisterDeleteObject(mcps, writeClient, algoliaIndexName)
	records.RegisterGetObject(mcps, writeClient, algoliaIndexName)
	records.RegisterInsertObject(mcps, writeClient, algoliaIndexName)
	records.RegisterInsertObjects(mcps, writeClient, algoliaIndexName)

	// Tools for searching
	query.RegisterRunQuery(mcps, client, algoliaIndexName)

	// Tools for managing rules
	rules.RegisterDeleteRule(mcps, writeClient, algoliaIndexName)
	rules.RegisterSearchRules(mcps, client, algoliaIndexName)

	// Tools for managing synonyms
	synonyms.RegisterClearSynonyms(mcps, writeClient, algoliaIndexName)
	synonyms.RegisterDeleteSynonym(mcps, writeClient, algoliaIndexName)
	synonyms.RegisterGetSynonym(mcps, client, algoliaIndexName)
	synonyms.RegisterInsertSynonym(mcps, writeClient, algoliaIndexName)
	synonyms.RegisterSearchSynonym(mcps, client, algoliaIndexName)

	if err := server.ServeStdio(mcps); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
