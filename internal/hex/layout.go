package hex

import (
	hx "github.com/gojuno/go.hexgrid"
	morton "github.com/gojuno/go.morton"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IsoScaleX = 1.0
	IsoScaleY = 0.5
)

var (
	OffsetX = 640.0
	OffsetY = 360.0
	HexSize = 40.0
)

type Layout struct {
	hexgrid *hx.Grid
}

func NewLayout() *Layout {
	center := hx.MakePoint(0, 0)
	size := hx.MakePoint(HexSize, HexSize)
	grid := hx.MakeGrid(hx.OrientationFlat, center, size, morton.Make64(2, 32))
	return &Layout{
		hexgrid: grid,
	}
}

func (l *Layout) HexToPixel(q, r int64) (float64, float64) {
	currentGrid := hx.MakeGrid(hx.OrientationFlat, hx.MakePoint(0, 0), hx.MakePoint(HexSize, HexSize), morton.Make64(2, 32))
	hex := hx.MakeHex(q, r)

	point := currentGrid.HexCenter(hex)

	x := point.X()*IsoScaleX + OffsetX
	y := point.Y()*IsoScaleY + OffsetY

	return x, y
}

func (l *Layout) PixelToHex(x, y float64) (int64, int64) {
	currentGrid := hx.MakeGrid(hx.OrientationFlat, hx.MakePoint(0, 0), hx.MakePoint(HexSize, HexSize), morton.Make64(2, 32))
	x -= OffsetX
	y -= OffsetY

	x /= IsoScaleX
	y /= IsoScaleY

	point := hx.MakePoint(x, y)
	hex := currentGrid.HexAt(point)

	return hex.Q(), hex.R()
}

func (l *Layout) GetCorners(q, r int64) []ebiten.Vertex {
	currentGrid := hx.MakeGrid(hx.OrientationFlat, hx.MakePoint(0, 0), hx.MakePoint(HexSize, HexSize), morton.Make64(2, 32))
	hex := hx.MakeHex(q, r)
	corners := currentGrid.HexCorners(hex)

	vertices := make([]ebiten.Vertex, 6)

	for i, corner := range corners {
		x := corner.X()*IsoScaleX + OffsetX
		y := corner.Y()*IsoScaleY + OffsetY

		vertices[i] = ebiten.Vertex{
			DstX:   float32(x),
			DstY:   float32(y),
			SrcX:   0, // Will be used when we add textures
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		}

	}

	return vertices
}

func (l *Layout) ZoomIn() {
	HexSize *= 1.1
}

func (l *Layout) ZoomOut() {
	HexSize /= 1.1
}

func GetVisibleHexes(origin Hex, maxRange int64, isBlocking func(Hex) bool) []Hex {
	visibleHexes := []Hex{}

	for q := origin.Q - maxRange; q <= origin.Q+maxRange; q++ {
		for r := origin.R - maxRange; r <= origin.R+maxRange; r++ {
			target := Hex{Q: q, R: r}
			if HexDistance(origin, target) > maxRange {
				continue
			}

			if HasLineOfSight(origin, target, isBlocking) {
				visibleHexes = append(visibleHexes, target)
			}
		}
	}

	return visibleHexes
}

func HasLineOfSight(from, to Hex, isBlocking func(Hex) bool) bool {
	line := HexLine(from, to)

	// Check each hex along the line
	for _, h := range line {
		if h == to {
			return true // reached target
		}
		if isBlocking(h) {
			return false // blocked before reaching target
		}
	}
	return true
}

func HexLine(from, to Hex) []Hex {
	distance := HexDistance(from, to)
	results := []Hex{}

	for i := 0; i <= int(distance); i++ {
		t := float64(i) / float64(distance)
		q := from.Q + int64(float64(to.Q-from.Q)*t)
		r := from.R + int64(float64(to.R-from.R)*t)
		results = append(results, Hex{Q: q, R: r})
	}

	return results
}
