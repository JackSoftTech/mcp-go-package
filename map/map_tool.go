package mapTools

import (
	"context"
	"fmt"
	"sort"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MapParams struct {
	MapType   string         `json:"map_type"` // "square" or "hex"
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	TileTypes map[string]int `json:"tile_types"` // tile type -> count, use "any" for unlimited
}

var mapTypeMapping = map[string]func(MapResult) Map{
	(&SquareMap{}).Key(): GetSquareMap,
	(&HexMap{}).Key():    GetHexMap,
}

func ValidateMapType(mapType string) bool {
	_, ok := mapTypeMapping[mapType]
	return ok
}

func GenerateMapTool(ctx context.Context, req *mcp.CallToolRequest, args MapParams) (*mcp.CallToolResult, any, error) {
	// Validate map type
	if !ValidateMapType(args.MapType) {
		args.MapType = (&SquareMap{}).Key() // default type
	}

	gameMap, err := GenerateRandomMap(args)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: " + err.Error()},
			},
		}, nil, err
	}
	var result string
	if info, ok := interface{}(gameMap).(interface{ Info() string }); ok {
		result += info.Info() + "\n"
	}
	tiles := gameMap.GetTiles()
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Y == tiles[j].Y {
			return tiles[i].X < tiles[j].X
		}
		return tiles[i].Y < tiles[j].Y
	})
	for _, tile := range tiles {
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
