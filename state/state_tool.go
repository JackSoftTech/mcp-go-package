package state

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigureStateParams struct {
	Schema  map[string]any `json:"schema"`
	Initial map[string]any `json:"initial,omitempty"`
}

type SetStateParams struct {
	State map[string]any `json:"state"`
}

type GetStateParams struct{}

func ConfigureStateTool(ctx context.Context, req *mcp.CallToolRequest, args ConfigureStateParams) (*mcp.CallToolResult, any, error) {
	parsedSchema, err := ParseStateSchema(args.Schema)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil, err
	}
	if err := ConfigureStateSchema(parsedSchema, args.Initial); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "state schema configured"},
		},
	}, nil, nil
}

func SetStateTool(ctx context.Context, req *mcp.CallToolRequest, args SetStateParams) (*mcp.CallToolResult, any, error) {
	if err := SetDynamicState(args.State); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "state updated"},
		},
	}, nil, nil
}

func GetCurrentStateTool(ctx context.Context, req *mcp.CallToolRequest, args GetStateParams) (*mcp.CallToolResult, any, error) {
	state, err := GetDynamicState()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil, err
	}
	payload, err := json.Marshal(state)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(payload)},
		},
	}, nil, nil
}

func AddConfigureStateToolHandler(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "configure_state",
		Description: "configure state schema and initial values",
	}, ConfigureStateTool)
}

func AddSetStateToolHandler(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_state",
		Description: "set current state with schema validation",
	}, SetStateTool)
}

func AddGetCurrentStateToolHandler(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_state",
		Description: "get current state as json",
	}, GetCurrentStateTool)
}
