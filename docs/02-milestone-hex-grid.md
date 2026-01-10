# Milestone 2: Hex Grid Rendering

## Goal

Create an isometric hex grid that you can see and interact with using the mouse.

## Learning Objectives

- Understand hex coordinate systems (axial coordinates)
- Convert between screen pixels and hex coordinates
- Render hexagons in an isometric view
- Handle mouse input to select hexes
- Create a reusable hex grid system

## Prerequisites

- Completed Milestone 1
- Installed `pmcxs/hexgrid` package

## Tasks

### Hex Grid Foundation
- [ ] Create `internal/hex/layout.go` for coordinate conversion
- [ ] Define hex size and orientation constants
- [ ] Implement pixel-to-hex conversion function
- [ ] Implement hex-to-pixel conversion function
- [ ] Write tests for coordinate conversion (see `internal/hex/hex_test.go`)

### Grid Rendering
- [ ] Create a hex texture/sprite (or draw programmatically)
- [ ] Render a 10x10 hex grid
- [ ] Apply isometric transformation
- [ ] Color hexes in a checkerboard pattern
- [ ] Add coordinate labels to hexes

### Mouse Interaction
- [ ] Detect mouse position
- [ ] Convert mouse position to hex coordinates
- [ ] Highlight hex under mouse cursor
- [ ] Click to select a hex
- [ ] Show selected hex coordinates

### Testing
- [ ] Run tests: `go test ./internal/hex`
- [ ] Grid renders correctly
- [ ] Mouse hover highlights correct hex
- [ ] Clicking selects the right hex
- [ ] Coordinate labels match expected positions

## Hex Grid Theory

### Why Hexagons?

Hexagons are superior to squares for games because:
- Each neighbor is equidistant (perfect for movement)
- 6 neighbors instead of 4 (or 8 with diagonals)
- No "diagonal" movement issues
- Look natural and organic

### Coordinate Systems

We'll use **axial coordinates** (q, r):

```
     r
    ↗
   / \
  /   \
 q→   (0,0)
```

- `q`: column (moves right)
- `r`: row (moves down-right)
- Alternative: cube coordinates (x, y, z where x+y+z=0)

Example grid:
```
  (-1,-1) (0,-1) (1,-1)
    (-1,0) (0,0) (1,0)
  (-1,1) (0,1) (1,1)
```

### Flat-Top vs Pointy-Top

Hexagons have two orientations:

**Pointy-Top** (we'll use this):
```
   / \
  |   |
   \ /
```

**Flat-Top**:
```
  ____
 /    \
 \____/
```

Pointy-top works better for isometric view.

### Isometric Projection

Normal hex grid is top-down. Isometric tilts it to look 3D:

```
Top-down:        Isometric:
  / \              /\
 |   |    →       /  \
  \ /            /    \
```

We achieve this by:
1. Rotating the coordinate system
2. Squashing the Y-axis (typically by 0.5)

## Step-by-Step Implementation

### Step 1: Install hexgrid Library

```bash
go get github.com/pmcxs/hexgrid
```

### Step 2: Create Hex Layout Package

Create `internal/hex/layout.go`:

```go
package hex

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pmcxs/hexgrid"
)

const (
	// Hex size in pixels (radius from center to corner)
	HexSize = 40.0

	// Isometric scaling factor
	IsoScaleX = 1.0
	IsoScaleY = 0.5 // Squash Y axis for isometric look

	// Offset for centering the grid on screen
	OffsetX = 640.0
	OffsetY = 200.0
)

// Layout handles conversion between hex and pixel coordinates
type Layout struct {
	hexgrid *hexgrid.Grid
}

// NewLayout creates a new hex layout with pointy-top orientation
func NewLayout() *Layout {
	// Create hexgrid with pointy-top orientation
	grid := hexgrid.NewGrid(hexgrid.FlatOrientation, HexSize)

	return &Layout{
		hexgrid: grid,
	}
}

// HexToPixel converts hex coordinates (q, r) to screen pixels
func (l *Layout) HexToPixel(q, r int) (float64, float64) {
	// Use hexgrid to get base position
	hex := hexgrid.NewHex(q, r)
	point := l.hexgrid.Center(hex)

	// Apply isometric transformation
	x := point.X * IsoScaleX
	y := point.Y * IsoScaleY

	// Add screen offset to center the grid
	x += OffsetX
	y += OffsetY

	return x, y
}

// PixelToHex converts screen pixels to hex coordinates (q, r)
func (l *Layout) PixelToHex(x, y float64) (int, int) {
	// Remove screen offset
	x -= OffsetX
	y -= OffsetY

	// Reverse isometric transformation
	x /= IsoScaleX
	y /= IsoScaleY

	// Use hexgrid to convert to hex coordinates
	point := hexgrid.NewPoint(x, y)
	hex := l.hexgrid.PixelToHex(point)

	// Round to nearest hex (hexgrid returns fractional coordinates)
	return hex.Q, hex.R
}

// GetCorners returns the 6 corner points of a hex in pixel coordinates
func (l *Layout) GetCorners(q, r int) []ebiten.Vertex {
	hex := hexgrid.NewHex(q, r)
	corners := l.hexgrid.Corners(hex)

	vertices := make([]ebiten.Vertex, 6)
	for i, corner := range corners {
		// Apply isometric transformation
		x := corner.X*IsoScaleX + OffsetX
		y := corner.Y*IsoScaleY + OffsetY

		vertices[i] = ebiten.Vertex{
			DstX: float32(x),
			DstY: float32(y),
			SrcX: 0, // We'll set these when we use textures
			SrcY: 0,
			ColorR: 1.0,
			ColorG: 1.0,
			ColorB: 1.0,
			ColorA: 1.0,
		}
	}

	return vertices
}
```

### Step 3: Create Test File

Create `internal/hex/hex_test.go`:

```go
package hex

import (
	"math"
	"testing"
)

func TestHexToPixel(t *testing.T) {
	layout := NewLayout()

	tests := []struct {
		name string
		q, r int
		// We can't test exact pixels due to isometric transform,
		// but we can verify the function doesn't panic
	}{
		{"origin", 0, 0},
		{"positive q", 1, 0},
		{"positive r", 0, 1},
		{"negative q", -1, 0},
		{"negative r", 0, -1},
		{"diagonal", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := layout.HexToPixel(tt.q, tt.r)

			// Basic sanity checks
			if math.IsNaN(x) || math.IsNaN(y) {
				t.Errorf("HexToPixel(%d, %d) returned NaN", tt.q, tt.r)
			}
			if math.IsInf(x, 0) || math.IsInf(y, 0) {
				t.Errorf("HexToPixel(%d, %d) returned Inf", tt.q, tt.r)
			}
		})
	}
}

func TestPixelToHex(t *testing.T) {
	layout := NewLayout()

	tests := []struct {
		name string
		x, y float64
		// Expected hex coordinates
		wantQ, wantR int
	}{
		{"center of screen", OffsetX, OffsetY, 0, 0},
		// Add more tests after you verify the first one works
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r := layout.PixelToHex(tt.x, tt.y)

			if q != tt.wantQ || r != tt.wantR {
				t.Errorf("PixelToHex(%f, %f) = (%d, %d), want (%d, %d)",
					tt.x, tt.y, q, r, tt.wantQ, tt.wantR)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	layout := NewLayout()

	// Converting hex -> pixel -> hex should give us back the original hex
	tests := []struct {
		q, r int
	}{
		{0, 0},
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
		{3, 3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Hex to pixel
			x, y := layout.HexToPixel(tt.q, tt.r)

			// Pixel back to hex
			q, r := layout.PixelToHex(x, y)

			if q != tt.q || r != tt.r {
				t.Errorf("Round trip failed: (%d,%d) -> (%f,%f) -> (%d,%d)",
					tt.q, tt.r, x, y, q, r)
			}
		})
	}
}
```

Run the tests:
```bash
go test ./internal/hex
```

### Step 4: Update Game to Render Grid

Update `cmd/game/main.go`:

```go
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/yourusername/boss-battler/internal/hex" // Update with your module path
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridSize     = 5 // 5x5 hex grid
)

type Game struct {
	layout       *hex.Layout
	hoveredQ     int
	hoveredR     int
	selectedQ    int
	selectedR    int
	hasSelection bool
}

func NewGame() *Game {
	return &Game{
		layout:    hex.NewLayout(),
		selectedQ: -999, // Invalid hex
		selectedR: -999,
	}
}

func (g *Game) Update() error {
	// Quit on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Get mouse position
	mx, my := ebiten.CursorPosition()

	// Convert to hex coordinates
	g.hoveredQ, g.hoveredR = g.layout.PixelToHex(float64(mx), float64(my))

	// Click to select
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Check if we clicked on a valid grid hex
		if g.isValidHex(g.hoveredQ, g.hoveredR) {
			g.selectedQ = g.hoveredQ
			g.selectedR = g.hoveredR
			g.hasSelection = true
		}
	}

	return nil
}

func (g *Game) isValidHex(q, r int) bool {
	// Check if hex is within our grid bounds
	return q >= -gridSize && q <= gridSize && r >= -gridSize && r <= gridSize
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw all hexes in the grid
	for q := -gridSize; q <= gridSize; q++ {
		for r := -gridSize; r <= gridSize; r++ {
			g.drawHex(screen, q, r)
		}
	}

	// Draw UI
	msg := fmt.Sprintf("Milestone 2: Hex Grid\nHovered: (%d, %d)\n", g.hoveredQ, g.hoveredR)
	if g.hasSelection {
		msg += fmt.Sprintf("Selected: (%d, %d)\n", g.selectedQ, g.selectedR)
	} else {
		msg += "Click to select a hex\n"
	}
	msg += "\nPress ESC to quit"

	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) drawHex(screen *ebiten.Image, q, r int) {
	// Get center position
	cx, cy := g.layout.HexToPixel(q, r)

	// Choose color based on state
	var hexColor color.Color

	if g.hasSelection && g.selectedQ == q && g.selectedR == r {
		// Selected hex - bright green
		hexColor = color.RGBA{100, 255, 100, 255}
	} else if g.hoveredQ == q && g.hoveredR == r {
		// Hovered hex - yellow
		hexColor = color.RGBA{255, 255, 100, 255}
	} else {
		// Checkerboard pattern for visual interest
		if (q+r)%2 == 0 {
			hexColor = color.RGBA{60, 60, 80, 255}
		} else {
			hexColor = color.RGBA{50, 50, 70, 255}
		}
	}

	// Draw hexagon using vector graphics
	// For simplicity, draw as circle (you can draw proper hex with corners later)
	vector.DrawFilledCircle(screen, float32(cx), float32(cy), hex.HexSize*0.8, hexColor, false)

	// Draw outline
	vector.StrokeCircle(screen, float32(cx), float32(cy), hex.HexSize*0.8, 2, color.RGBA{100, 100, 120, 255}, false)

	// Draw coordinates
	coordText := fmt.Sprintf("%d,%d", q, r)
	// Note: This text positioning is approximate
	ebitenutil.DebugPrintAt(screen, coordText, int(cx)-15, int(cy)-5)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boss Battler - Hex Grid")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

### Step 5: Run and Test

```bash
go run cmd/game/main.go
```

You should see:
- A grid of hexagons (drawn as circles for now)
- Hexes in a checkerboard pattern
- Hovering highlights hexes in yellow
- Clicking selects a hex (green)
- Coordinates displayed on each hex

### Step 6: Draw Proper Hexagons (Optional Enhancement)

The circles work, but let's draw actual hexagons. Update the `drawHex` function:

```go
func (g *Game) drawHex(screen *ebiten.Image, q, r int) {
	// Get corner points
	corners := g.layout.GetCorners(q, r)

	// Choose color
	var hexColor color.Color
	if g.hasSelection && g.selectedQ == q && g.selectedR == r {
		hexColor = color.RGBA{100, 255, 100, 255}
	} else if g.hoveredQ == q && g.hoveredR == r {
		hexColor = color.RGBA{255, 255, 100, 255}
	} else {
		if (q+r)%2 == 0 {
			hexColor = color.RGBA{60, 60, 80, 255}
		} else {
			hexColor = color.RGBA{50, 50, 70, 255}
		}
	}

	// Draw filled hexagon using vector
	// Convert corners to path
	var path vector.Path

	// Start at first corner
	path.MoveTo(corners[0].DstX, corners[0].DstY)

	// Draw lines to other corners
	for i := 1; i < 6; i++ {
		path.LineTo(corners[i].DstX, corners[i].DstY)
	}

	// Close the path
	path.Close()

	// Fill the hexagon
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	// Set color for all vertices
	r, g, b, a := hexColor.RGBA()
	for i := range vertices {
		vertices[i].ColorR = float32(r) / 0xffff
		vertices[i].ColorG = float32(g) / 0xffff
		vertices[i].ColorB = float32(b) / 0xffff
		vertices[i].ColorA = float32(a) / 0xffff
	}

	screen.DrawTriangles(vertices, indices, emptyImage, nil)

	// Draw outline
	outlineColor := color.RGBA{100, 100, 120, 255}
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen,
			corners[i].DstX, corners[i].DstY,
			corners[next].DstX, corners[next].DstY,
			2, outlineColor, false)
	}

	// Draw coordinates
	cx, cy := g.layout.HexToPixel(q, r)
	coordText := fmt.Sprintf("%d,%d", q, r)
	ebitenutil.DebugPrintAt(screen, coordText, int(cx)-15, int(cy)-5)
}

// Add this helper at package level
var emptyImage = ebiten.NewImage(3, 3)

func init() {
	emptyImage.Fill(color.White)
}
```

This will draw proper hexagons instead of circles.

## Key Concepts

### Axial Coordinates

Axial coords (q, r) are like tilted x, y:
- `q` moves right
- `r` moves down-right
- Each hex has 6 neighbors at (q±1, r±0), (q±0, r±1), (q∓1, r±1)

### Coordinate Conversion

The hexgrid library handles most of the math, but understanding helps:

**Hex to Pixel** (simplified):
```
x = size * (√3 * q + √3/2 * r)
y = size * (3/2 * r)
```

**Pixel to Hex**:
- More complex (involves rounding)
- hexgrid.PixelToHex handles this

### Isometric Transformation

We apply after basic hex math:
```go
isoX = x * scaleX
isoY = y * scaleY // scaleY < 1 for isometric "tilt"
```

### Distance Between Hexes

```go
func hexDistance(q1, r1, q2, r2 int) int {
	// Convert to cube coordinates
	x1, y1, z1 := q1, -q1-r1, r1
	x2, y2, z2 := q2, -q2-r2, r2

	// Manhattan distance in cube space / 2
	return (abs(x1-x2) + abs(y1-y2) + abs(z1-z2)) / 2
}
```

## Common Issues

### Hexes don't align with mouse
- Check your isometric transform is consistent
- Verify OffsetX/OffsetY values
- Test with the round-trip test (hex->pixel->hex)

### Grid looks stretched
- Adjust IsoScaleY (typically 0.5 for isometric)
- Try different values: 0.4, 0.6, 0.7

### Wrong hex highlights
- PixelToHex might need tweaking
- Check coordinate system (pointy vs flat top)
- Verify hexgrid orientation setting

### Performance issues
- Drawing many vector shapes can be slow
- For large grids, only draw visible hexes
- Consider using pre-rendered hex sprites

## Next Steps

Milestone 2 complete! You now have:
- A working hex grid system
- Coordinate conversion
- Mouse interaction
- Foundation for positioning game entities

In [Milestone 3](03-milestone-ecs-setup.md), you'll add the Entity Component System to start placing characters and enemies on the hex grid!

## Extra Challenges (Optional)

- [ ] Highlight hexes within range of selected hex
- [ ] Implement hex pathfinding (A*)
- [ ] Draw hex grid with sprite textures
- [ ] Add zoom in/out functionality
- [ ] Show distance from selected hex to hovered hex
- [ ] Implement hex field-of-view algorithm
