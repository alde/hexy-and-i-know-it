package main

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	c "github.com/alde/hexy-and-i-know-it/internal/color"
	"github.com/alde/hexy-and-i-know-it/internal/hex"
)

const (
	screenWidth        = 1280
	screenHeight       = 720
	gridSize     int64 = 5
)

var (
	ebitenImage *ebiten.Image
	emptyImage  *ebiten.Image
	hexImage    *ebiten.Image
)

func init() {
	emptyImage = ebiten.NewImage(1, 1)
	emptyImage.Fill(color.White)

	stoneImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			// Add some noise/variation to base grey
			val := 50 + rand.Intn(30)
			stoneImg.Set(x, y, color.RGBA{uint8(val), uint8(val), uint8(val), 255})
		}
	}
	ebitenImage = ebiten.NewImageFromImage(stoneImg)
	hexImage = ebitenImage
}

type Game struct {
	updateCount int
	bgColor     color.Color
	debug       bool

	layout               *hex.Layout
	hoveredQ, hoveredR   int64
	selectedQ, selectedR int64

	hasSelection               bool
	pathFromSelectionToHovered []hex.Hex
	visibleHexes               []hex.Hex
}

func NewGame() *Game {
	return &Game{
		bgColor:                    color.RGBA{30, 30, 40, 255},
		layout:                     hex.NewLayout(),
		selectedQ:                  -999,
		selectedR:                  -999,
		pathFromSelectionToHovered: []hex.Hex{},
	}
}

func (g *Game) Update() error {
	g.updateCount++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	mx, my := ebiten.CursorPosition()
	g.hoveredQ, g.hoveredR = g.layout.PixelToHex(float64(mx), float64(my))
	if g.hasSelection {
		g.pathFromSelectionToHovered = hex.FindPath(
			hex.Hex{Q: g.selectedQ, R: g.selectedR},
			hex.Hex{Q: g.hoveredQ, R: g.hoveredR},
			func(h hex.Hex) bool {
				return g.isValidHex(h.Q, h.R)
			},
		)
		g.visibleHexes = hex.GetVisibleHexes(
			hex.Hex{Q: g.selectedQ, R: g.selectedR},
			3,
			func(h hex.Hex) bool {
				return false // No blocking hexes for now
			},
		)

	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.selectHex(g.hoveredQ, g.hoveredR)
	}

	if ebiten.IsKeyPressed(ebiten.KeyAlt) {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.debug = !g.debug
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if hexImage == emptyImage {
			hexImage = ebitenImage
		} else {
			hexImage = emptyImage
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		g.layout.ZoomIn()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		g.layout.ZoomOut()
	}

	return nil
}

func (g *Game) selectHex(q, r int64) {
	if g.isValidHex(g.hoveredQ, g.hoveredR) {
		g.selectedQ = g.hoveredQ
		g.selectedR = g.hoveredR
		g.hasSelection = true
	}
}

func (g *Game) isValidHex(q, r int64) bool {
	return q >= -gridSize && q <= gridSize && r >= -gridSize && r <= gridSize
}

func isAdjacent(q1, r1, q2, r2 int64) bool {
	dq := q1 - q2
	dr := r1 - r2
	if (dq == 1 && dr == 0) || (dq == -1 && dr == 0) ||
		(dq == 0 && dr == 1) || (dq == 0 && dr == -1) ||
		(dq == 1 && dr == -1) || (dq == -1 && dr == 1) {
		return true
	}
	return false
}

func hexListContains(hexes []hex.Hex, q, r int64) bool {
	for _, h := range hexes {
		if h.Q == q && h.R == r {
			return true
		}
	}
	return false
}

func (g *Game) drawHex(screen *ebiten.Image, q, r int64) {
	corners := g.layout.GetCorners(q, r)

	adjacentHexes := []hex.Hex{}
	if g.hasSelection {
		adjacentHexes = hex.GetVisibleHexes(
			hex.Hex{Q: g.selectedQ, R: g.selectedR}, 1, func(h hex.Hex) bool { return false },
		)
	}

	var hexColor color.Color
	if g.hasSelection && q == g.selectedQ && r == g.selectedR {
		hexColor = c.Color0
	} else if g.hoveredQ == q && g.hoveredR == r {
		hexColor = c.Color3
	} else if hexListContains(g.pathFromSelectionToHovered, q, r) {
		hexColor = c.Color4
	} else if hexListContains(adjacentHexes, q, r) {
		hexColor = c.Color1
	} else if hexListContains(g.visibleHexes, q, r) {
		hexColor = c.Color7
	} else {
		if (q+r)%2 == 0 {
			hexColor = c.Color5
		} else {
			hexColor = c.Color6
		}
	}
	var path vector.Path
	path.MoveTo(corners[0].DstX, corners[0].DstY)
	for i := 1; i < len(corners); i++ {
		path.LineTo(corners[i].DstX, corners[i].DstY)
	}
	path.Close()
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	red, green, blue, alpha := hexColor.RGBA()
	bounds := hexImage.Bounds()
	for i := range vertices {
		vertices[i].ColorR = float32(red) / 0xffff
		vertices[i].ColorG = float32(green) / 0xffff
		vertices[i].ColorB = float32(blue) / 0xffff
		vertices[i].ColorA = float32(alpha) / 0xffff

		vertices[i].SrcX = float32(bounds.Min.X + i*bounds.Dx()/len(vertices))
		vertices[i].SrcY = float32(bounds.Min.Y + i*bounds.Dy()/len(vertices))
	}

	screen.DrawTriangles(vertices, indices, hexImage, nil)

	outlineColor := color.White
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, corners[i].DstX, corners[i].DstY, corners[next].DstX, corners[next].DstY, 2, outlineColor, false)
	}

	if g.debug {
		cx, cy := g.layout.HexToPixel(q, r)
		coordText := fmt.Sprintf("(%d,%d)", q, r)
		// Approximate text positioning
		ebitenutil.DebugPrintAt(screen, coordText, int(cx)-15, int(cy)-5)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	for q := -gridSize; q <= gridSize; q++ {
		for r := -gridSize; r <= gridSize; r++ {
			g.drawHex(screen, q, r)
		}
	}

	msg := fmt.Sprintf("Milestone 2 - Hex Grid\nHovered Hex: (%d, %d)", g.hoveredQ, g.hoveredR)
	if g.hasSelection {
		msg += fmt.Sprintf("\nSelected Hex: (%d, %d)", g.selectedQ, g.selectedR)
		msg += fmt.Sprintf("\nDistance from selected to hovered: %d", len(g.pathFromSelectionToHovered)-1)
	} else {
		msg += "\nClick a hex to select it"
	}
	msg += "\nPress ALT+D to toggle debug info\nPress ALT+ENTER to toggle fullscreen\nPress ESC to quit"

	ebitenutil.DebugPrintAt(screen, msg, 10, 10)

	if g.debug {
		mouseX, mouseY := ebiten.CursorPosition()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: %d, %d", mouseX, mouseY), 10, screenHeight-80)
		msg := fmt.Sprintf("FPS: %0.2f\nTPS: %0.2f\nUpdates: %d", ebiten.ActualFPS(), ebiten.ActualTPS(), g.updateCount)
		ebitenutil.DebugPrintAt(screen, msg, 10, screenHeight-60)
	}

	ebiten.SetWindowTitle(fmt.Sprintf("Hexy and I know it. %0.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, ousideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hexy and I Know It")

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("failed to run game", "error", err)
		os.Exit(1)
	}
}
