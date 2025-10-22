package mapTools

import "fmt"

type HexMap struct {
	Width  int    // Number of hexagonal tiles horizontally
	Height int    // Number of hexagonal tiles vertically
	Tiles  []Tile // One-dimensional array of tiles with coordinates
}

func (m *HexMap) Key() string {
	return "hex"
}

// Info returns the type and dimensions of the map
func (m *HexMap) Info() string {
	return fmt.Sprintf("HexMap: width=%d, height=%d", m.Width, m.Height)
}

func (m *HexMap) GetTiles() []Tile {
	return m.Tiles
}

func GetHexMap(res MapResult) Map {
	return &HexMap{
		Width:  res.opts.Width,
		Height: res.opts.Height,
		Tiles:  res.Tiles,
	}
}
