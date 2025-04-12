package connectors

import (
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/ingestion"
	"github.com/google/uuid"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/algolia/mcp/pkg/mcputil"
)

func RegisterListConnectors(mcps *server.MCPServer, client *ingestion.APIClient) {
	getObjectTool := mcp.NewTool(
		"list_connector",
		mcp.WithDescription("List all existing connectors"),
	)

	mcps.AddTool(getObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sources, err := client.ListSources(client.NewApiListSourcesRequest())
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not list sources: %v", err),
			), nil
		}

		return mcputil.JSONToolResult("connectors", sources)
	})
}

func RegisterTaskForAConnector(mcps *server.MCPServer, client *ingestion.APIClient) {
	getObjectTool := mcp.NewTool(
		"list_tasks_for_a_connector",
		mcp.WithDescription("List all tasks for a source or a connector"),
		mcp.WithString(
			"source_id",
			mcp.Description("The source or connector id to look up"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceID, _ := req.Params.Arguments["source_id"].(string)

		listTasks, err := client.ListTasks(client.NewApiListTasksRequest())
		if err != nil {
			return nil, err
		}

		if len(listTasks.Tasks) == 0 {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not find any tasks"),
			), nil
		}

		var tasks []string
		for _, task := range listTasks.Tasks {
			if task.SourceID == sourceID {
				tasks = append(tasks, task.TaskID)
			}
		}

		return mcputil.JSONToolResult("tasks_for_this_source", tasks)
	})
}

func RegisterStartTask(mcps *server.MCPServer, client *ingestion.APIClient) {
	getObjectTool := mcp.NewTool(
		"start_task",
		mcp.WithDescription("Start a task"),
		mcp.WithString(
			"task_id",
			mcp.Description("The task id to start"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		taskID, _ := req.Params.Arguments["task_id"].(string)

		task, err := client.RunTask(client.NewApiRunTaskRequest(taskID))
		if err != nil {
			return nil, err
		}

		return mcputil.JSONToolResult("task", task)
	})
}

func RegisterCreateNewConnector(mcps *server.MCPServer, defaultIndex string, apiKey string, appid string, client *ingestion.APIClient) {
	getObjectTool := mcp.NewTool(
		"create_new_json_connector",
		mcp.WithDescription("Create a new JSON connector"),
		mcp.WithString(
			"url",
			mcp.Description("The URL to the hosted json file"),
			mcp.Required(),
		),
		mcp.WithString(
			"unique_column_name",
			mcp.Description("The json property we will use as unique identifier (Algolia ObjectID afterward)"),
			mcp.Required(),
		),
		mcp.WithString(
			"index_name",
			mcp.Description("Algolia index name you want to use or create."),
			mcp.Required(),
			mcp.DefaultString(defaultIndex),
		),
		mcp.WithString("scheduled_cron",
			mcp.Description("The cron schedule to run the task"),
		),
	)

	mcps.AddTool(getObjectTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := req.Params.Arguments["url"].(string)
		uniqueColumnID := req.Params.Arguments["unique_column_name"].(string)
		indexName, _ := req.Params.Arguments["index_name"].(string)

		method := ingestion.METHOD_TYPE_GET

		connector, err := client.CreateSource(client.NewApiCreateSourceRequest(&ingestion.SourceCreate{
			Type: "json",
			Name: "Test Connector from LLM " + uuid.New().String(),
			Input: &ingestion.SourceInput{
				SourceJSON: &ingestion.SourceJSON{
					Url:            url,
					UniqueIDColumn: &uniqueColumnID,
					Method:         &method,
				},
			},
			AuthenticationID: nil,
		}))

		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("unable to create source: %v", err),
			), nil
		}

		auth, err := client.CreateAuthentication(client.NewApiCreateAuthenticationRequest(&ingestion.AuthenticationCreate{
			Type: "algolia",
			Name: "Algolia Search API Key" + uuid.New().String(),
			Input: ingestion.AuthInput{
				AuthAlgolia: &ingestion.AuthAlgolia{
					AppID:  appid,
					ApiKey: apiKey,
				},
			},
		},
		))

		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("unable to create authentication: %v", err),
			), nil
		}

		destination, err := client.CreateDestination(client.NewApiCreateDestinationRequest(&ingestion.DestinationCreate{
			Type: "search",
			Name: "Index from Claude" + uuid.New().String(),
			Input: ingestion.DestinationInput{
				DestinationIndexName: &ingestion.DestinationIndexName{
					IndexName: indexName,
				},
			},
			AuthenticationID:  &auth.AuthenticationID,
			TransformationIDs: nil,
		}))

		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("unable to create your destination: %v", err),
			), nil
		}

		i, ok := req.Params.Arguments["scheduled_cron"]
		var cron *string
		if ok {
			s := i.(string)
			cron = &s
		}
		t := true
		task, err := client.CreateTask(client.NewApiCreateTaskRequest(&ingestion.TaskCreate{
			SourceID:      connector.SourceID,
			DestinationID: destination.DestinationID,
			Action:        ingestion.ACTION_TYPE_REPLACE,
			Cron:          cron,
			Enabled:       &t,
		}))

		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("unable to create task: %v", err),
			), nil
		}

		mcputil.JSONResource(connector)
		mcputil.JSONResource(destination)
		mcputil.JSONResource(task)

		return mcputil.JSONToolResult("task", task)
	})
}
