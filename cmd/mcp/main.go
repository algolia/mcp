package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/abtesting"
	"github.com/algolia/mcp/pkg/analytics"
	"github.com/algolia/mcp/pkg/collections"
	"github.com/algolia/mcp/pkg/monitoring"
	"github.com/algolia/mcp/pkg/querysuggestions"
	"github.com/algolia/mcp/pkg/recommend"
	searchpkg "github.com/algolia/mcp/pkg/search"
	"github.com/algolia/mcp/pkg/usage"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create a new MCP server with name and version
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "Algolia MCP",
		Version: "0.0.2",
	}, nil)

	// Parse MCP_ENABLED_TOOLS environment variable to determine which toolsets to enable
	enabledToolsEnv := os.Getenv("MCP_ENABLED_TOOLS")
	enabled := make(map[string]bool)
	allTools := []string{"abtesting", "analytics", "collections", "monitoring", "querysuggestions", "recommend", "search", "search_read", "search_write", "usage"}

	// If MCP_ENABLED_TOOLS is set, enable only the specified toolsets
	// Otherwise, enable all toolsets
	if enabledToolsEnv == "" {
		for _, toolName := range allTools {
			enabled[toolName] = true
		}
	}

	for _, toolName := range strings.Split(enabledToolsEnv, ",") {
		trimmedName := strings.ToLower(strings.TrimSpace(toolName))
		for _, knownTool := range allTools {
			if trimmedName == knownTool {
				enabled[trimmedName] = true
				break
			}
		}
	}

	// Initialize Algolia client
	var searchClient *search.Client
	var searchIndex *search.Index

	// Get Algolia credentials from environment variables
	appID := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")
	indexName := os.Getenv("ALGOLIA_INDEX_NAME")

	searchClient = search.NewClient(appID, apiKey)
	searchIndex = searchClient.InitIndex(indexName)

	// Register tools from enabled packages.
	if enabled["abtesting"] {
		abtesting.RegisterTools(mcpServer)
	}
	if enabled["analytics"] {
		analytics.RegisterTools(mcpServer)
	}
	if enabled["collections"] {
		collections.RegisterTools(mcpServer)
	}
	if enabled["monitoring"] {
		monitoring.RegisterTools(mcpServer)
	}
	if enabled["querysuggestions"] {
		querysuggestions.RegisterAll(mcpServer)
	}
	if enabled["recommend"] {
		recommend.RegisterAll(mcpServer)
	}
	if enabled["search"] {
		searchpkg.RegisterAll(mcpServer)
	} else {
		// Only register specific search tools if "search" is not enabled
		if enabled["search_read"] {
			searchpkg.RegisterReadAll(mcpServer, searchClient, searchIndex)
		}
		if enabled["search_write"] {
			searchpkg.RegisterWriteAll(mcpServer, searchClient, searchIndex)
		}
	}
	if enabled["usage"] {
		usage.RegisterAll(mcpServer)
	}

	// Create a logger that writes to stderr instead of stdout
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// Log to stderr to avoid interfering with JSON-RPC communication
	logger.Println("Starting MCP server...")

	// Check server type from environment variable (defaults to "stdio" if not set)
	serverType := strings.ToLower(strings.TrimSpace(os.Getenv("MCP_SERVER_TYPE")))

	// The official SDK primarily supports stdio transport
	// For now, we'll use stdio transport as the main method
	if serverType != "" && serverType != "stdio" {
		logger.Printf("Warning: Server type '%s' not fully supported with official SDK, defaulting to stdio", serverType)
	}

	// Log to stderr to avoid interfering with JSON-RPC communication
	logger.Println("Starting stdio server...")

	// Use stdio transport with the official SDK
	if err := mcpServer.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		logger.Fatalf("MCP server failed: %v", err)
	}
}
