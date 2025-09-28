package mapTools

import "fmt"

type SquareMap struct {
	Width  int    // Number of square tiles horizontally
	Height int    // Number of square tiles vertically
	Tiles  []Tile // One-dimensional array of tiles with coordinates
}

// Info returns the type and dimensions of the map
func (m *SquareMap) Info() string {
	return fmt.Sprintf("SquareMap: width=%d, height=%d", m.Width, m.Height)
}

func (m *SquareMap) GetTiles() []Tile {
	return m.Tiles
}

func (m *SquareMap) GetType() MapType {
	return Square
}
