package victoryTools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EndStateParams struct {
	EndStateType string `json:"end_state_type"`
	ScoreTarget  int    `json:"score_target"`
}

var endStateTypeMapping = map[string]func(EndStateResult) EndState{
	(&ScoreEndState{}).Key(): SetScoreEndState,
}

func ValidateEndStateType(endStateType string) bool {
	_, ok := endStateTypeMapping[endStateType]
	return ok
}

func SetupEndStateTool(ctx context.Context, req *mcp.CallToolRequest, args EndStateParams) (*mcp.CallToolResult, any, error) {
	if !ValidateEndStateType(args.EndStateType) {
		args.EndStateType = (&ScoreEndState{}).Key()
	}

	endState := endStateTypeMapping[args.EndStateType](EndStateResult{opts: args})
	info := endState.Info()
	result := fmt.Sprintf("%s", info)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func AddSetupEndStateToolHandler(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "setup_end_state",
		Description: "setup an end state",
	}, SetupEndStateTool)
}
