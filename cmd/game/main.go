package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Game struct {
	updateCount int
	bgColor     color.Color
}

func (g *Game) Update() error {
	g.updateCount++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	rect := ebiten.NewImage(200, 100)
	rect.Fill(color.RGBA{100, 150, 200, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(50, 50)
	screen.DrawImage(rect, opts)

	ebitenutil.DebugPrint(screen, "Hexy and I Know It\nPress SPACE to change background color\nPress ESC to quit")

	msg := fmt.Sprintf("FPS: %0.2f, TPS: %0.2f, Updates: %d", ebiten.ActualFPS(), ebiten.ActualTPS(), g.updateCount)

	ebitenutil.DebugPrintAt(screen, msg, 10, screenHeight-60)
}

func (g *Game) Layout(outsideWidth, ousideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		bgColor: color.RGBA{30, 30, 40, 255},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hexy and I Know It")

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("failed to run game", "error", err)
		os.Exit(1)
	}
}
