package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/alde/hexy-and-i-know-it/internal/hex"
)

const (
	screenWidth        = 1280
	screenHeight       = 720
	gridSize     int64 = 5
)

var emptyImage = ebiten.NewImage(3, 3)

func init() {
	emptyImage.Fill(color.White)
}

type Game struct {
	updateCount int
	bgColor     color.Color
	debug       bool

	Hex struct {
		layout               *hex.Layout
		hoveredQ, hoveredR   int64
		selectedQ, selectedR int64
	}

	hasSelection bool
}

func NewGame() *Game {
	return &Game{
		bgColor: color.RGBA{30, 30, 40, 255},

		Hex: struct {
			layout               *hex.Layout
			hoveredQ, hoveredR   int64
			selectedQ, selectedR int64
		}{
			layout:    hex.NewLayout(),
			selectedQ: -999,
			selectedR: -999,
		},
	}
}

func (g *Game) Update() error {
	g.updateCount++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	mx, my := ebiten.CursorPosition()
	g.Hex.hoveredQ, g.Hex.hoveredR = g.Hex.layout.PixelToHex(float64(mx), float64(my))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.isValidHex(g.Hex.hoveredQ, g.Hex.hoveredR) {
			g.Hex.selectedQ = g.Hex.hoveredQ
			g.Hex.selectedR = g.Hex.hoveredR
			g.hasSelection = true
		}
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
		if g.bgColor == (color.RGBA{30, 30, 40, 255}) {
			g.bgColor = color.RGBA{40, 30, 50, 255}
		} else {
			g.bgColor = color.RGBA{30, 30, 40, 255}
		}
	}

	return nil
}

func (g *Game) isValidHex(q, r int64) bool {
	return q >= -gridSize && q <= gridSize && r >= -gridSize && r <= gridSize
}

func (g *Game) drawHex(screen *ebiten.Image, q, r int64) {
	corners := g.Hex.layout.GetCorners(q, r)

	var hexColor color.Color

	if g.hasSelection && g.Hex.selectedQ == q && g.Hex.selectedR == r {
		hexColor = color.RGBA{100, 255, 100, 255}
	} else if g.Hex.hoveredQ == q && g.Hex.hoveredR == r {
		hexColor = color.RGBA{255, 255, 100, 255}
	} else {
		if (q+r)%2 == 0 {
			hexColor = color.RGBA{60, 60, 80, 255}
		} else {
			hexColor = color.RGBA{50, 50, 70, 255}
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
	for i := range vertices {
		vertices[i].ColorR = float32(red) / 0xffff
		vertices[i].ColorG = float32(green) / 0xffff
		vertices[i].ColorB = float32(blue) / 0xffff
		vertices[i].ColorA = float32(alpha) / 0xffff
	}
	screen.DrawTriangles(vertices, indices, emptyImage, nil)

	outlineColor := color.RGBA{100, 100, 120, 255}
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, corners[i].DstX, corners[i].DstY, corners[next].DstX, corners[next].DstY, 2, outlineColor, false)
	}

	cx, cy := g.Hex.layout.HexToPixel(q, r)
	coordText := fmt.Sprintf("(%d,%d)", q, r)
	// Approximate text positioning
	ebitenutil.DebugPrintAt(screen, coordText, int(cx)-15, int(cy)-5)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	for q := -gridSize; q <= gridSize; q++ {
		for r := -gridSize; r <= gridSize; r++ {
			g.drawHex(screen, q, r)
		}
	}

	msg := fmt.Sprintf("Milestone 2 - Hex Grid\nHovered Hex: (%d, %d)", g.Hex.hoveredQ, g.Hex.hoveredR)
	if g.hasSelection {
		msg += fmt.Sprintf("\nSelected Hex: (%d, %d)", g.Hex.selectedQ, g.Hex.selectedR)
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
