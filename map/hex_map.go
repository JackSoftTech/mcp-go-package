package mapTools

type HexMap struct {
	Width  int    // Number of hexagonal tiles horizontally
	Height int    // Number of hexagonal tiles vertically
	Tiles  []Tile // One-dimensional array of tiles with coordinates
}
