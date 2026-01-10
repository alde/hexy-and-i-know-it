# Milestone 1: Hello Ebiten

## Goal

Get Ebiten running with a basic game window and understand the core game loop pattern.

## Learning Objectives

- Understand Ebiten's `Game` interface
- Learn the Update/Draw separation
- See how TPS (ticks per second) and FPS (frames per second) work
- Handle basic keyboard input
- Display text and shapes

## Tasks

### Setup
- [ ] Create `cmd/game/main.go` file
- [ ] Implement minimal `Game` struct
- [ ] Open a 1280x720 window with title "Boss Battler"
- [ ] Run the game and verify window opens

### Game Loop
- [ ] Implement `Update()` method (game logic, runs 60 times/second)
- [ ] Implement `Draw()` method (rendering only)
- [ ] Implement `Layout()` method (define logical screen size)
- [ ] Add frame counter to track updates

### Basic Rendering
- [ ] Draw a colored background
- [ ] Display "Boss Battler" text on screen
- [ ] Show FPS/TPS counter
- [ ] Draw a simple rectangle

### Input Handling
- [ ] Detect ESC key to quit game
- [ ] Detect SPACE key to change background color
- [ ] Print input events to console

### Testing
- [ ] Window opens without errors
- [ ] Game runs at 60 TPS
- [ ] ESC key quits the game
- [ ] SPACE key changes background
- [ ] No console errors

## Step-by-Step Implementation

### Step 1: Create Main File

Create `cmd/game/main.go`:

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

// Game implements the ebiten.Game interface
type Game struct {
	// We'll add fields here as we go
}

// Update is called every tick (60 times per second)
// This is where game logic goes
func (g *Game) Update() error {
	// TODO: Implement game logic
	return nil
}

// Draw is called every frame
// This is where rendering goes
func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: Implement rendering
}

// Layout returns the logical screen size
// Ebiten will scale the screen to fit the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boss Battler")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

**Run it**:
```bash
go run cmd/game/main.go
```

You should see a black window. Press Ctrl+C to quit for now.

### Step 2: Add State to Game Struct

Add fields to track game state:

```go
type Game struct {
	updateCount int
	bgColor     color.Color
}
```

Initialize in main:

```go
func main() {
	game := &Game{
		bgColor: color.RGBA{30, 30, 40, 255}, // Dark blue-gray
	}
	// ... rest of main
}
```

### Step 3: Implement Update Logic

```go
func (g *Game) Update() error {
	g.updateCount++

	// Quit on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Change background on SPACE
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Toggle between two colors
		if g.bgColor == (color.RGBA{30, 30, 40, 255}) {
			g.bgColor = color.RGBA{40, 30, 50, 255} // Purple-ish
		} else {
			g.bgColor = color.RGBA{30, 30, 40, 255} // Back to blue-gray
		}
	}

	return nil
}
```

Don't forget to import `inpututil`:
```go
import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)
```

### Step 4: Implement Draw Method

```go
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background
	screen.Fill(g.bgColor)

	// Draw a rectangle (we'll use this for UI elements later)
	rect := ebiten.NewImage(200, 100)
	rect.Fill(color.RGBA{100, 150, 200, 200}) // Semi-transparent blue
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(50, 50) // Position at (50, 50)
	screen.DrawImage(rect, opts)

	// Draw text
	ebitenutil.DebugPrint(screen, "Boss Battler - Milestone 1\nPress SPACE to change color\nPress ESC to quit")

	// Show FPS/TPS
	msg := fmt.Sprintf("FPS: %0.2f, TPS: %0.2f\nUpdates: %d",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		g.updateCount)
	ebitenutil.DebugPrintAt(screen, msg, 10, screenHeight-60)
}
```

Add imports:
```go
import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)
```

### Step 5: Test Your Game

Run it:
```bash
go run cmd/game/main.go
```

Verify:
- Window opens with title "Boss Battler"
- Background is dark blue-gray
- Blue rectangle appears in top-left
- Text shows controls and FPS/TPS
- SPACE key changes background color
- ESC key quits the game
- FPS should be around 60

## Complete Code Reference

<details>
<summary>Click to see complete main.go</summary>

```go
package main

import (
	"fmt"
	"image/color"
	"log"

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

	// Quit on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Change background on SPACE
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
	// Fill background
	screen.Fill(g.bgColor)

	// Draw a rectangle
	rect := ebiten.NewImage(200, 100)
	rect.Fill(color.RGBA{100, 150, 200, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(50, 50)
	screen.DrawImage(rect, opts)

	// Draw text
	ebitenutil.DebugPrint(screen, "Boss Battler - Milestone 1\nPress SPACE to change color\nPress ESC to quit")

	// Show FPS/TPS
	msg := fmt.Sprintf("FPS: %0.2f, TPS: %0.2f\nUpdates: %d",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		g.updateCount)
	ebitenutil.DebugPrintAt(screen, msg, 10, screenHeight-60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		bgColor: color.RGBA{30, 30, 40, 255},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boss Battler")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

</details>

## Key Concepts

### The Game Loop

Ebiten runs two methods repeatedly:

1. **Update()** - Called 60 times per second (60 TPS)
   - Handle input
   - Update game state
   - Run physics, AI, combat calculations
   - Modify data
   - Return `ebiten.Termination` to quit

2. **Draw()** - Called as fast as possible (usually 60 FPS)
   - Read game state
   - Render to screen
   - **Never modify game state here**

### Why Separate Update and Draw?

- **Update** runs at fixed rate (deterministic, good for physics/logic)
- **Draw** runs as fast as possible (smooth visuals)
- On slow devices, Draw might skip frames but Update always runs 60x/sec
- This prevents "slow motion" game logic on low-end hardware

### TPS vs FPS

- **TPS** (Ticks Per Second): How often Update() runs (always 60)
- **FPS** (Frames Per Second): How often Draw() runs (varies by hardware)
- If FPS < 60, you might have performance issues
- If TPS < 60, your Update() is taking too long (bad!)

### Input Handling

Two ways to check input:

1. **IsKeyPressed** - True while key is held
   ```go
   if ebiten.IsKeyPressed(ebiten.KeyEscape) {
       // Runs every frame while ESC is held
   }
   ```

2. **IsKeyJustPressed** - True only on first frame of press
   ```go
   if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
       // Runs once per key press (good for toggles)
   }
   ```

### Performance Tips

- **Don't create images in Draw()** - Very slow, causes stuttering
- **Create images once** (in initialization or Update), reuse in Draw()
- Bad example:
  ```go
  func (g *Game) Draw(screen *ebiten.Image) {
      rect := ebiten.NewImage(100, 100) // SLOW: Creates image every frame
      screen.DrawImage(rect, nil)
  }
  ```
- Good example:
  ```go
  type Game struct {
      rect *ebiten.Image
  }

  func NewGame() *Game {
      return &Game{
          rect: ebiten.NewImage(100, 100), // Create once
      }
  }

  func (g *Game) Draw(screen *ebiten.Image) {
      screen.DrawImage(g.rect, nil) // Reuse existing image
  }
  ```

## Common Issues

### Window doesn't open
- Check you have all dependencies: `go get github.com/hajimehoshi/ebiten/v2`
- Make sure you're in the right directory
- Check for compile errors: `go build cmd/game/main.go`

### FPS is low
- Check `ebiten.ActualTPS()` - should be close to 60
- If TPS is low, your Update() is too slow
- Profile with `go run -tags=ebitenginedebug cmd/game/main.go`

### Can't quit with ESC
- Make sure you're returning `ebiten.Termination` not `nil`
- Check the key name is `ebiten.KeyEscape`

## Next Steps

Milestone 1 complete! You now understand:
- Ebiten's game loop structure
- Update vs Draw separation
- Basic rendering and input
- TPS/FPS concepts

In [Milestone 2](02-milestone-hex-grid.md), you'll render an isometric hex grid and handle mouse clicks to select hexes. This is where the game really starts to take shape!

## Extra Challenges (Optional)

Want to explore more? Try these:

- [ ] Add mouse position tracking and display it
- [ ] Draw multiple rectangles at different positions
- [ ] Create a simple animation (moving rectangle)
- [ ] Add more keyboard shortcuts
- [ ] Change window title to show FPS
- [ ] Try fullscreen mode: `ebiten.SetFullscreen(true)`
