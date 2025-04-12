package main

import (
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/ingestion"
	"github.com/algolia/mcp/pkg/connectors"
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

	algoliaAppRegion := os.Getenv("ALGOLIA_APP_REGION")

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
	indices.RegisterGetSettings(mcps, index)

	// Tools for managing records
	records.RegisterGetObject(mcps, index)
	records.RegisterInsertObject(mcps, writeIndex)
	records.RegisterInsertObjects(mcps, writeIndex)

	// Tools for searching
	query.RegisterRunQuery(mcps, client, index)

	// Tools for managing rules
	rules.RegisterSearchRules(mcps, index)

	// Tools for managing synonyms
	synonyms.RegisterSearchSynonym(mcps, index)

	// Connectors
	newClient, err := ingestion.NewClient(algoliaAppID, algoliaWriteAPIKey, ingestion.Region(algoliaAppRegion))
	if err != nil {
		log.Fatalf("Error creating ingestion client: %v", err)
	}

	connectors.RegisterListConnectors(mcps, newClient)
	connectors.RegisterTaskForAConnector(mcps, newClient)
	connectors.RegisterStartTask(mcps, newClient)
	connectors.RegisterCreateNewConnector(mcps, algoliaIndexName, algoliaWriteAPIKey, algoliaAppID, newClient)

	if err := server.ServeStdio(mcps); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
