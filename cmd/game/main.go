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
	debug       bool
}

type ImageAndOpts struct {
	image *ebiten.Image
	opts  *ebiten.DrawImageOptions
}

var imageStore = map[string]*ImageAndOpts{}

func (g *Game) Update() error {
	g.updateCount++

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	for name, imgAndOpts := range imageStore {
		if name == "exampleRect1" {

			imgAndOpts.opts.GeoM.Reset()
			imgAndOpts.opts.GeoM.Translate(
				75+float64(g.updateCount%200),
				75+float64(g.updateCount%150),
			)
		}
		screen.DrawImage(imgAndOpts.image, imgAndOpts.opts)
	}

	ebitenutil.DebugPrint(screen, "Hexy and I Know It\nPress SPACE to change background color\nPress ESC to quit")

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
	game := &Game{
		bgColor: color.RGBA{30, 30, 40, 255},
	}

	imageStore = map[string]*ImageAndOpts{
		"exampleRect0": createRectangle(200, 100, color.RGBA{100, 150, 200, 200}, 50, 50),
		"exampleRect1": createRectangle(300, 200, color.RGBA{200, 100, 150, 200}, 75, 75),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hexy and I Know It")

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("failed to run game", "error", err)
		os.Exit(1)
	}
}

func createRectangle(posX, posY int, color color.Color, sizeX, sizeY int) *ImageAndOpts {
	rect := ebiten.NewImage(posX, posY)
	rect.Fill(color)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(sizeX), float64(sizeY))
	return &ImageAndOpts{
		image: rect,
		opts:  opts,
	}
}
