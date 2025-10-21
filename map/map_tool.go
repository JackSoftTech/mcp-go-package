package mapTools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GenerateMapParams struct {
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	MapType   string   `json:"map_type"`   // "square" or "hex"
	TileTypes []string `json:"tile_types"` // tile type list
}

func GenerateMapTool(ctx context.Context, req *mcp.CallToolRequest, args GenerateMapParams) (*mcp.CallToolResult, any, error) {
	var mapType MapType
	switch args.MapType {
	case "hex":
		mapType = Hex
	default:
		mapType = Square
	}
	opts := MapGenOptions{
		MapType:   mapType,
		Width:     args.Width,
		Height:    args.Height,
		TileTypes: args.TileTypes,
	}
	gameMap := GenerateRandomMap(opts)
	var result string
	if info, ok := interface{}(gameMap).(interface{ Info() string }); ok {
		result += info.Info() + "\n"
	}
	for _, tile := range gameMap.GetTiles() {
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
