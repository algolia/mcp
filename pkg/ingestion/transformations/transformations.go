package transformations

import (
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/ingestion"
	"github.com/algolia/mcp/pkg/mcputil"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterListTransformations(mcps *server.MCPServer, client *ingestion.APIClient) {
	listTransformationsTool := mcp.NewTool(
		"list_transformations",
		mcp.WithDescription("List all existing JS transformations"),
	)

	mcps.AddTool(listTransformationsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		transformations, err := client.ListTransformations(client.NewApiListTransformationsRequest())
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not list transformations: %v", err),
			), nil
		}

		return mcputil.JSONToolResult("transformations", transformations)
	})
}

func RegisterGetTransformation(mcps *server.MCPServer, client *ingestion.APIClient) {
	getTransformationTool := mcp.NewTool(
		"get_transformation",
		mcp.WithDescription("Get a transformation by its ID"),
		mcp.WithString(
			"transformation_id",
			mcp.Description("The transformation id to look up"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getTransformationTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		transformationID, _ := req.Params.Arguments["transformation_id"].(string)

		transformation, err := client.GetTransformation(client.NewApiGetTransformationRequest(transformationID))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not get transformation: %v", err),
			), nil
		}

		return mcputil.JSONToolResult("transformation", transformation)
	})
}

func RegisterUpdateTransformation(mcps *server.MCPServer, client *ingestion.APIClient) {
	getTransformationTool := mcp.NewTool(
		"update_transformation",
		mcp.WithDescription("Update a transformation by its ID. You need to take all the parameters of the previously fetched transformation, not just the ones you want to update"),
		mcp.WithString(
			"transformation_id",
			mcp.Description("The transformation id to look up"),
			mcp.Required(),
		),
		mcp.WithString(
			"code",
			mcp.Description("The code of the transformation"),
			mcp.Required(),
		),
		mcp.WithString(
			"name",
			mcp.Description("The name of the transformation"),
			mcp.Required(),
		),
		mcp.WithString(
			"description",
			mcp.Description("The description of the transformation"),
		),
	)

	mcps.AddTool(getTransformationTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		transformationID, _ := req.Params.Arguments["transformation_id"].(string)
		code, _ := req.Params.Arguments["code"].(string)
		name, _ := req.Params.Arguments["name"].(string)
		description, _ := req.Params.Arguments["description"].(string)
		authenticationIDs, _ := req.Params.Arguments["authentication_ids"].([]string)

		transformation, err := client.UpdateTransformation(client.NewApiUpdateTransformationRequest(transformationID, &ingestion.TransformationCreate{
			Code:              code,
			Name:              name,
			Description:       &description,
			AuthenticationIDs: authenticationIDs,
		}))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not update transformation: %v", err),
			), nil
		}

		return mcputil.JSONToolResult("transformation", transformation)
	})
}

func RegisterTryTransformation(mcps *server.MCPServer, client *ingestion.APIClient) {
	getTransformationTool := mcp.NewTool(
		"try_transformation",
		mcp.WithDescription("Execute the JavaScript transformation code with a specific sample dataset."),
		mcp.WithString(
			"code",
			mcp.Description("The javascript code we want to try"),
			mcp.Required(),
		),
		mcp.WithObject(
			"sample",
			mcp.Description("the JSON object we want to use as a sample"),
			mcp.Required(),
		),
	)

	mcps.AddTool(getTransformationTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		code, _ := req.Params.Arguments["code"].(string)
		sampleRecord, _ := req.Params.Arguments["sample"].(map[string]interface{})

		transformation, err := client.TryTransformation(client.NewApiTryTransformationRequest(&ingestion.TransformationTry{
			Code:            code,
			SampleRecord:    sampleRecord,
			Authentications: nil,
		}))
		if err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("could not update transformation: %v", err),
			), nil
		}

		return mcputil.JSONToolResult("transformed sample", transformation)
	})
}
