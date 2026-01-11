package hex

import (
	hx "github.com/gojuno/go.hexgrid"
	morton "github.com/gojuno/go.morton"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IsoScaleX = 1.0
	IsoScaleY = 0.5

	OffsetX = 640.0
	OffsetY = 360.0
)

const HexSize = 40.0

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
	hex := hx.MakeHex(q, r)
	point := l.hexgrid.HexCenter(hex)

	x := point.X()*IsoScaleX + OffsetX
	y := point.Y()*IsoScaleY + OffsetY

	return x, y
}

func (l *Layout) PixelToHex(x, y float64) (int64, int64) {
	x -= OffsetX
	y -= OffsetY

	x /= IsoScaleX
	y /= IsoScaleY

	point := hx.MakePoint(x, y)
	hex := l.hexgrid.HexAt(point)

	return hex.Q(), hex.R()
}

func (l *Layout) GetCorners(q, r int64) []ebiten.Vertex {
	hex := hx.MakeHex(q, r)
	corners := l.hexgrid.HexCorners(hex)

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
