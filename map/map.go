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

type MapType int

const (
	Square MapType = iota
	Hex
)

type Map interface {
	Info() string
	GetTiles() []Tile
	GetType() MapType
}

func GenerateRandomMap(mapType MapType, width, height int) Map {
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
	if mapType == Hex {
		return &HexMap{
			Width:  width,
			Height: height,
			Tiles:  tiles,
		}
	}
	return &SquareMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}
