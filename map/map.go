package mapTools

import (
	"fmt"
	"math/rand"
	"time"
)

type MapResult struct {
	opts  MapParams
	Tiles []Tile
}

type Map interface {
	Key() string
	Info() string
	GetTiles() []Tile
}

func GenerateRandomMap(opts MapParams) (Map, error) {
	rand.Seed(time.Now().UnixNano())
	totalTiles := opts.Width * opts.Height
	fixedCount := 0
	var flexibleTypes []string
	for k, count := range opts.TileTypes {
		if count >= 0 {
			fixedCount += count
		} else {
			flexibleTypes = append(flexibleTypes, k)
		}
	}
	if fixedCount > totalTiles {
		return nil, fmt.Errorf("tile count (%d) exceeds map size (%d)", fixedCount, totalTiles)
	}
	remaining := totalTiles - fixedCount
	// 隨機分配剩餘格子給 flexibleTypes
	flexibleCounts := make(map[string]int)
	for _, k := range flexibleTypes {
		flexibleCounts[k] = 0
	}
	for i := 0; i < remaining; i++ {
		t := flexibleTypes[rand.Intn(len(flexibleTypes))]
		flexibleCounts[t]++
	}
	coords := make([][2]int, 0, totalTiles)
	for x := 0; x < opts.Width; x++ {
		for y := 0; y < opts.Height; y++ {
			coords = append(coords, [2]int{x, y})
		}
	}
	// shuffle coords
	for i := range coords {
		j := rand.Intn(i + 1)
		coords[i], coords[j] = coords[j], coords[i]
	}
	tiles := make([]Tile, 0, totalTiles)
	idx := 0
	// 固定數量的 tile type
	for tileType, count := range opts.TileTypes {
		if count < 0 {
			continue
		}
		for c := 0; c < count && idx < len(coords); c++ {
			x, y := coords[idx][0], coords[idx][1]
			tiles = append(tiles, Tile{
				Type: tileType,
				X:    x,
				Y:    y,
			})
			idx++
		}
	}
	// 隨機分配的 tile type
	for tileType, count := range flexibleCounts {
		for c := 0; c < count && idx < len(coords); c++ {
			x, y := coords[idx][0], coords[idx][1]
			tiles = append(tiles, Tile{
				Type: tileType,
				X:    x,
				Y:    y,
			})
			idx++
		}
	}
	if idx != totalTiles {
		return nil, fmt.Errorf("tile count does not match map size")
	}
	mapResult := mapTypeMapping[opts.MapType](MapResult{
		opts:  opts,
		Tiles: tiles,
	})

	return mapResult, nil
}
