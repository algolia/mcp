# MCP

This repository contains experimental [Model Context Protocol (or MCP)](https://modelcontextprotocol.io/introduction) servers for interacting with Algolia APIs.

## Installation

First follow the [quick start](https://modelcontextprotocol.io/quickstart/user), which will install Claude Desktop and setup a sample Fileserver MCP server.  This is a great introduction to using MCP and will let you debug things using the official guide if there are issues.

## Setup the prototype Algolia MCP server

Requirements:

* Go (https://go.dev/doc/install)

### Clone the repo and build the server

Clone the repo, amd build the mcp server:

```shell
$ git clone git@github.com:algolia/mcp.git
$ cd mcp/cmd/mcp
$ go build
```
We need to have the full path of the built server:
```shell
$ pwd
/path/to/the/repo/cmd/mcp
```
The full path to the built MCP server is:

```shell
/path/to/the/repo/cmd/mcp/mcp
```

### Update the settings to point to the new server

In Claude desktop edit the settings as per https://modelcontextprotocol.io/quickstart/user#2-add-the-filesystem-mcp-server and this time add the server definition for algolia (using the server path that you found earlier).

```json
{
   "mcpServers": {
      "algolia": {
         "command": "/path/to/the/repo/cmd/mcp/mcp",
         "env": {
            "ALGOLIA_APP_ID": "<APP_ID>",
            "ALGOLIA_INDEX_NAME": "<INDEX_NAME>",
            "ALGOLIA_API_KEY": "<API_KEY>",
            "ALGOLIA_WRITE_API_KEY": "<ADMIN_API_KEY>"  /* if you want to allow write operations, use your ADMIN key here */
         }
      }
   }
}
```

Restart Claude desktop, and you should see a new `"algolia"` tool is available.

## Debugging

You can run the Inspector (see https://modelcontextprotocol.io/docs/tools/inspector) to check the MCP features and run them manually.

From the repo root, setup the environment

```shell
$ export ALGOLIA_APP_ID=""
$ export ALGOLIA_INDEX_NAME=""
$ export ALGOLIA_API_KEY=""
$ export ALGOLIA_WRITE_API_KEY=""  # if you want to allow write operations, use your ADMIN key here
```
Move into the server directory, and rebuild (if necessary):
```shell
$ cd cmd/mcp
$ go build # might already be up-to-date
```
Run the MCP inspector on the server:
```shell
$ npx @modelcontextprotocol/inspector ./mcp
```

## Using with Ollama

You can actually run a local mcphost (which orchestrates the MCP servers for you), and then use them with other models locally via Ollama.

We are using https://github.com/mark3labs/mcphost for this.

As per the [README](https://github.com/mark3labs/mcphost?tab=readme-ov-file#installation-) you need a a config file, so you can copy the Claude one, and put it somewhere sensible so you can use it on the command line (for example `~/mcp.json`)

```json filename="~/mcp.json"
{
   "mcpServers": {
      "algolia": {
         "command": "/path/to/the/repo/cmd/mcp/mcp",
         "env": {
            "ALGOLIA_APP_ID": "<APP_ID>",
            "ALGOLIA_INDEX_NAME": "<INDEX_NAME>",
            "ALGOLIA_API_KEY": "<API_KEY>"
         }
      }
   }
}
```
You can now run it directly (no need to check out the repo):
```shell
$ go run github.com/mark3labs/mcphost@latest --config ~/mcp.json -m ollama:qwen2.5:3b
```

# FAQ

* Resource templates and root are not supported by Claude desktop right now: https://github.com/orgs/modelcontextprotocol/discussions/136
This is a weird one, since there is a bunch of content online showing the templates, maybe it's just not GA yet.
