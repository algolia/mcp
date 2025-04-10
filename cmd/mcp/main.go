package main

import (
	"fmt"
	"log"
	"os"

	amcp "github.com/algolia/mcp"

	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
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

	amcp.RegisterRunQuery(mcps, index)
	amcp.RegisterGetObject(mcps, index)
	amcp.RegisterGetSettings(mcps, index)
	amcp.RegisterSearchRules(mcps, index)
	amcp.RegisterSearchSynonym(mcps, index)

	amcp.RegisterInsertObject(mcps, writeIndex)
	amcp.RegisterInsertObjects(mcps, writeIndex)

	if err := server.ServeStdio(mcps); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
