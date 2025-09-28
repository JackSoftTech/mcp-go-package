package mapTools

import (
	"math/rand"
	"time"
)

const (
	Grass    = "grass"
	Forest   = "forest"
	Mountain = "mountain"
	Field    = "field"
	Desert   = "desert"
	Hill     = "hill"
)

var TileTypes = []string{Grass, Forest, Mountain, Field, Desert, Hill}

type Tile struct {
	Type string
	X    int
	Y    int
}

type Map struct {
	Tiles  []Tile
	Width  int
	Height int
}

func GenerateRandomMap(width, height int) Map {
	rand.Seed(time.Now().UnixNano())
	tiles := make([]Tile, 0, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			tileType := TileTypes[rand.Intn(len(TileTypes))]
			tiles = append(tiles, Tile{
				Type: tileType,
				X:    x,
				Y:    y,
			})
		}
	}
	return Map{
		Tiles:  tiles,
		Width:  width,
		Height: height,
	}
}
