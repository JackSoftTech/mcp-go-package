package mapTools

import "fmt"

type HexMap struct {
	Width  int    // Number of hexagonal tiles horizontally
	Height int    // Number of hexagonal tiles vertically
	Tiles  []Tile // One-dimensional array of tiles with coordinates
}

// Info returns the type and dimensions of the map
func (m *HexMap) Info() string {
	return fmt.Sprintf("HexMap: width=%d, height=%d", m.Width, m.Height)
}

func (m *HexMap) GetTiles() []Tile {
	return m.Tiles
}

func (m *HexMap) GetType() MapType {
	return Hex
}
