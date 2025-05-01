package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/mcp/pkg/abtesting"
	"github.com/algolia/mcp/pkg/analytics"
	"github.com/algolia/mcp/pkg/collections"
	"github.com/algolia/mcp/pkg/monitoring"
	"github.com/algolia/mcp/pkg/querysuggestions"
	"github.com/algolia/mcp/pkg/recommend"
	searchpkg "github.com/algolia/mcp/pkg/search"
	"github.com/algolia/mcp/pkg/usage"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server with name and version
	mcps := server.NewMCPServer("Algolia MCP", "0.0.2")

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
		abtesting.RegisterTools(mcps)
	}
	if enabled["analytics"] {
		analytics.RegisterTools(mcps)
	}
	if enabled["collections"] {
		collections.RegisterTools(mcps)
	}
	if enabled["monitoring"] {
		monitoring.RegisterTools(mcps)
	}
	if enabled["querysuggestions"] {
		querysuggestions.RegisterAll(mcps)
	}
	if enabled["recommend"] {
		recommend.RegisterAll(mcps)
	}
	if enabled["search"] {
		searchpkg.RegisterAll(mcps)
	} else {
		// Only register specific search tools if "search" is not enabled
		if enabled["search_read"] {
			searchpkg.RegisterReadAll(mcps, searchClient, searchIndex)
		}
		if enabled["search_write"] {
			searchpkg.RegisterWriteAll(mcps, searchClient, searchIndex)
		}
	}
	if enabled["usage"] {
		usage.RegisterAll(mcps)
	}

	// Create a logger that writes to stderr instead of stdout
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// Log to stderr to avoid interfering with JSON-RPC communication
	logger.Println("Starting MCP server...")

	// Check server type from environment variable (defaults to "stdio" if not set)
	serverType := strings.ToLower(strings.TrimSpace(os.Getenv("MCP_SERVER_TYPE")))

	// Start the appropriate server type
	if serverType == "sse" {
		// Get port from environment variable or use default
		portStr := os.Getenv("MCP_SSE_PORT")
		port := 8080 // Default port
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			} else {
				logger.Printf("Warning: Invalid MCP_SSE_PORT value '%s', using default port 8080", portStr)
			}
		}

		// Create the address string (e.g., ":8080")
		addr := fmt.Sprintf(":%d", port)
		logger.Printf("Starting SSE server on port %d...", port)

		// Create the SSE server
		sseServer := server.NewSSEServer(mcps)

		// Set up signal handling for graceful shutdown
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		// Start server in a goroutine
		serverErrCh := make(chan error, 1)
		go func() {
			if err := sseServer.Start(addr); err != nil && err != http.ErrServerClosed {
				serverErrCh <- fmt.Errorf("MCP server failed: %v", err)
				return
			}
			serverErrCh <- nil
		}()

		// Wait for either a shutdown signal or a server error
		select {
		case sig := <-signalChan:
			logger.Printf("Received signal %v, shutting down gracefully...", sig)
		case err := <-serverErrCh:
			if err != nil {
				logger.Fatalf("Server error: %v", err)
			}
		}

		// Use the server's shutdown method with a timeout context
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Attempt to shut down the server
		err := sseServer.Shutdown(shutdownCtx)

		// Always cancel the context to prevent resource leaks
		cancel()

		// Check for shutdown errors after ensuring context is canceled
		if err != nil {
			logger.Fatalf("Server shutdown failed: %v", err)
		}

		logger.Println("Server gracefully stopped")
	} else {
		// Default to stdio server
		if serverType != "" && serverType != "stdio" {
			logger.Printf("Warning: Unknown server type '%s', defaulting to stdio", serverType)
		}

		// Log to stderr to avoid interfering with JSON-RPC communication
		logger.Println("Starting stdio server...")

		// Use the same logger for error logging in the stdio server
		if err := server.ServeStdio(mcps, server.WithErrorLogger(logger)); err != nil {
			logger.Fatalf("MCP server failed: %v", err)
		}
	}
}
