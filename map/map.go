package mapTools

import (
	"math/rand"
	"time"
)

type MapType int

const (
	Square MapType = iota
	Hex
)

type MapGenOptions struct {
	MapType   MapType
	Width     int
	Height    int
	TileTypes []string
}

type Map interface {
	Info() string
	GetTiles() []Tile
	GetType() MapType
}

func GenerateRandomMap(opts MapGenOptions) Map {
	rand.Seed(time.Now().UnixNano())
	tiles := make([]Tile, 0, opts.Width*opts.Height)
	for x := 0; x < opts.Width; x++ {
		for y := 0; y < opts.Height; y++ {
			tileType := opts.TileTypes[rand.Intn(len(opts.TileTypes))]
			tiles = append(tiles, Tile{
				Type: tileType,
				X:    x,
				Y:    y,
			})
		}
	}
	if opts.MapType == Hex {
		return &HexMap{
			Width:  opts.Width,
			Height: opts.Height,
			Tiles:  tiles,
		}
	}
	return &SquareMap{
		Width:  opts.Width,
		Height: opts.Height,
		Tiles:  tiles,
	}
}
