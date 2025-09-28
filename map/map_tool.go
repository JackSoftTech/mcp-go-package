package mapTools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GenerateMapParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func GenerateMapTool(ctx context.Context, req *mcp.CallToolRequest, args GenerateMapParams) (*mcp.CallToolResult, any, error) {
	gameMap := GenerateRandomMap(args.Width, args.Height)
	var result string
	for _, tile := range gameMap.Tiles {
		result += fmt.Sprintf("Tile (%d, %d): %s\n", tile.X, tile.Y, tile.Type)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

// AddGenerateMapToolHandler 新增 generate_map Tool handler 到指定的 server
func AddGenerateMapToolHandler(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_map",
		Description: "generate a random map",
	}, GenerateMapTool)
}
