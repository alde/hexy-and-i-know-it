# Milestone 7: Demo Battle - Polish & Victory

## Goal

Create a complete, polished playable battle from start to finish with victory/defeat conditions, better AI, and visual improvements.

## Learning Objectives

- Implement victory/defeat detection
- Create simple boss AI behavior
- Add visual polish (animations, feedback)
- Improve UI and UX
- Handle game over states
- Make the game feel complete

## Prerequisites

- Completed Milestone 6
- All previous systems working

## Tasks

### Victory Conditions
- [ ] Detect when all enemies are defeated (victory)
- [ ] Detect when all players are defeated (defeat)
- [ ] Show victory/defeat screen
- [ ] Option to restart battle

### Boss AI Improvements
- [ ] Create smarter target selection (lowest HP, etc.)
- [ ] Add special boss abilities (optional)
- [ ] Make boss attack prioritize weakened targets

### Visual Polish
- [ ] Add damage number pop-ups
- [ ] Animate attacks (flash, shake)
- [ ] Show hit/miss indicators
- [ ] Improve entity sprites (draw actual hexagons instead of circles)
- [ ] Add color-coding for damage (red for hit, gray for miss)

### UI Improvements
- [ ] Better combat log formatting
- [ ] Show entity tooltips on hover
- [ ] Add HP bars above entities
- [ ] Improve turn indicator
- [ ] Add battle start screen

### Audio (Optional)
- [ ] Add hit sound effect
- [ ] Add miss sound effect
- [ ] Add victory/defeat music
- [ ] Background battle music

### Testing & Polish
- [ ] Play through full battle multiple times
- [ ] Fix any bugs found
- [ ] Balance damage/HP values
- [ ] Ensure game is fun to play

## Victory Detection Pattern

```
After each action:
    ↓
Check all player entities
    ↓
All dead? → DEFEAT
    ↓
Check all enemy entities
    ↓
All dead? → VICTORY
    ↓
Continue battle
```

## Step-by-Step Implementation

### Step 1: Add Victory Detection

Update `internal/states/battle.go`:

```go
// Add to BattleState enum:
const (
    BattleStart BattleState = iota
    RollInitiative
    PlayerTurn
    WaitingForPlayerInput
    EnemyTurn
    ExecuteAction
    Victory    // NEW
    Defeat     // NEW
    BattleEnd
)

// Update String():
func (s BattleState) String() string {
    return []string{
        "Battle Start",
        "Roll Initiative",
        "Player Turn",
        "Waiting For Input",
        "Enemy Turn",
        "Execute Action",
        "Victory!",
        "Defeat...",
        "Battle End",
    }[s]
}

// Add method to check victory conditions:
func (b *BattleManager) CheckVictoryConditions(world *donburi.World) {
    playerAlive := false
    enemyAlive := false

    // Check all combatants
    for _, entry := range b.TurnOrder {
        health := components.HealthComponent.Get(entry)

        if health.IsDead() {
            continue
        }

        if entry.HasComponent(components.PlayerControlledComponent) {
            playerAlive = true
        } else {
            enemyAlive = true
        }
    }

    if !playerAlive {
        b.State = Defeat
        b.AddLog("=== DEFEAT - All party members have fallen ===")
    } else if !enemyAlive {
        b.State = Victory
        b.AddLog("=== VICTORY - All enemies defeated! ===")
    }
}

// Update ExecutePendingAction to check victory after executing:
func (b *BattleManager) ExecutePendingAction(world *donburi.World) {
    if b.PendingAction == nil {
        return
    }

    result := b.PendingAction.Execute(world)

    for _, logEntry := range result.Logs {
        b.AddLog(logEntry)
    }

    b.PendingAction = nil

    // Check for victory/defeat
    b.CheckVictoryConditions(world)

    // Only advance turn if battle isn't over
    if b.State != Victory && b.State != Defeat {
        b.NextTurn()
    }
}
```

### Step 2: Improve Boss AI

Create `internal/systems/ai.go`:

```go
package systems

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/yourusername/boss-battler/internal/commands"
	"github.com/yourusername/boss-battler/internal/components"
)

// SelectBossTarget picks the best target for the boss to attack
func SelectBossTarget(world *donburi.World) *donburi.Entry {
	// Find all living player-controlled entities
	query := donburi.NewQuery(
		filter.Contains(
			components.PlayerControlledComponent,
			components.HealthComponent,
		),
	)

	var targets []*donburi.Entry
	query.Each(world, func(entry *donburi.Entry) {
		health := components.HealthComponent.Get(entry)
		if !health.IsDead() {
			targets = append(targets, entry)
		}
	})

	if len(targets) == 0 {
		return nil
	}

	// Strategy: Attack lowest HP target (finish them off)
	lowestHP := math.MaxInt
	var bestTarget *donburi.Entry

	for _, target := range targets {
		health := components.HealthComponent.Get(target)
		if health.Current < lowestHP {
			lowestHP = health.Current
			bestTarget = target
		}
	}

	return bestTarget
}

// CreateBossAction generates an action for the boss AI
func CreateBossAction(world *donburi.World, boss *donburi.Entry) commands.Action {
	target := SelectBossTarget(world)

	if target == nil {
		return &commands.WaitAction{Actor: boss}
	}

	return &commands.AttackAction{
		Attacker: boss,
		Target:   target,
	}
}
```

### Step 3: Add Damage Numbers

Update `internal/systems/render.go` to add a damage number system:

```go
// Add to package-level:
type DamageNumber struct {
	X, Y   float64
	Value  int
	TTL    int // Time to live in frames
	IsMiss bool
}

var damageNumbers = make([]*DamageNumber, 0)

// AddDamageNumber creates a floating damage number
func AddDamageNumber(x, y float64, damage int, isMiss bool) {
	damageNumbers = append(damageNumbers, &DamageNumber{
		X:      x,
		Y:      y,
		Value:  damage,
		TTL:    60, // 1 second at 60 FPS
		IsMiss: isMiss,
	})
}

// UpdateDamageNumbers should be called in Update()
func UpdateDamageNumbers() {
	for i := len(damageNumbers) - 1; i >= 0; i-- {
		dn := damageNumbers[i]
		dn.TTL--
		dn.Y -= 1 // Float upward

		if dn.TTL <= 0 {
			// Remove expired damage number
			damageNumbers = append(damageNumbers[:i], damageNumbers[i+1:]...)
		}
	}
}

// RenderDamageNumbers draws damage numbers
func RenderDamageNumbers(screen *ebiten.Image) {
	for _, dn := range damageNumbers {
		var text string
		var textColor color.Color

		if dn.IsMiss {
			text = "MISS"
			textColor = color.RGBA{150, 150, 150, 255} // Gray
		} else {
			text = fmt.Sprintf("-%d", dn.Value)
			textColor = color.RGBA{255, 100, 100, 255} // Red
		}

		// Fade out based on TTL
		alpha := uint8((float64(dn.TTL) / 60.0) * 255)
		if c, ok := textColor.(color.RGBA); ok {
			c.A = alpha
			textColor = c
		}

		// Draw text (note: ebitenutil.DebugPrintAt doesn't support colors,
		// so this is a simple version)
		ebitenutil.DebugPrintAt(screen, text, int(dn.X)-10, int(dn.Y))
	}
}
```

Update `internal/combat/attack.go` to trigger damage numbers:

```go
// At the end of PerformAttack, after applying damage:
// (You'll need to pass the hex layout to get screen coordinates)

// For now, we'll add this to the AttackResult so the caller can display it
type AttackResult struct {
    Hit          bool
    Critical     bool
    AttackRoll   int
    TotalAttack  int
    TargetAC     int
    Damage       int
    Killed       bool
    AttackerName string
    TargetName   string
    TargetEntry  *donburi.Entry // NEW: Add this so we can get position
}

// In PerformAttack, set:
result.TargetEntry = target
```

### Step 4: Add HP Bars

Update `internal/systems/render.go` to draw HP bars:

```go
// In RenderSystem, after drawing entity circle:

// Draw HP bar
if entry.HasComponent(components.HealthComponent) {
    health := components.HealthComponent.Get(entry)

    barWidth := float32(60)
    barHeight := float32(6)
    barX := float32(x) - barWidth/2
    barY := float32(y) + radius + 5

    // Background (red)
    vector.DrawFilledRect(screen, barX, barY, barWidth, barHeight,
        color.RGBA{100, 20, 20, 255}, false)

    // Foreground (green) based on HP percentage
    hpPercent := float32(health.Current) / float32(health.Max)
    if hpPercent < 0 {
        hpPercent = 0
    }

    // Color based on HP level
    var hpColor color.Color
    if hpPercent > 0.5 {
        hpColor = color.RGBA{50, 200, 50, 255} // Green
    } else if hpPercent > 0.25 {
        hpColor = color.RGBA{200, 200, 50, 255} // Yellow
    } else {
        hpColor = color.RGBA{200, 50, 50, 255} // Red
    }

    vector.DrawFilledRect(screen, barX, barY, barWidth*hpPercent, barHeight,
        hpColor, false)

    // Border
    vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 1,
        color.RGBA{255, 255, 255, 200}, false)
}
```

### Step 5: Update Main Game

Update `cmd/game/main.go`:

```go
func (g *Game) Update() error {
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        return ebiten.Termination
    }

    // Start battle on first frame
    if !g.started {
        g.battle.StartBattle(g.world)
        g.started = true
    }

    // Update damage numbers
    systems.UpdateDamageNumbers()

    // Handle different battle states
    switch g.battle.State {
    case states.WaitingForPlayerInput:
        activeCombatant := g.battle.GetActiveCombatant()
        if activeCombatant != nil {
            action := systems.ProcessPlayerInput(g.world, g.layout, g.inputState, activeCombatant)
            if action != nil {
                if err := g.battle.QueueAction(action); err != nil {
                    log.Printf("Invalid action: %v", err)
                }
            }
        }

    case states.ExecuteAction:
        g.battle.ExecutePendingAction(g.world)

        // Add damage number if there's a combat result
        // (This requires modifying how we pass attack results)

    case states.EnemyTurn:
        activeCombatant := g.battle.GetActiveCombatant()
        if activeCombatant != nil {
            action := systems.CreateBossAction(g.world, activeCombatant)
            g.battle.QueueAction(action)
        }

    case states.Victory, states.Defeat:
        // Press R to restart
        if inpututil.IsKeyJustPressed(ebiten.KeyR) {
            // Reset game
            *g = *NewGame()
            g.started = false
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{20, 20, 30, 255})

    // Draw hex grid
    for q := -gridSize; q <= gridSize; q++ {
        for r := -gridSize; r <= gridSize; r++ {
            g.drawHex(screen, q, r)
        }
    }

    // Draw entities with HP bars
    systems.RenderSystem(g.world, screen, g.layout)

    // Draw damage numbers
    systems.RenderDamageNumbers(screen)

    // Highlight hovered hex
    if g.inputState.HoveredHex != nil && g.battle.State == states.WaitingForPlayerInput {
        entity := systems.FindEntityAtHex(g.world, g.inputState.HoveredHex.Q, g.inputState.HoveredHex.R)
        if entity != nil {
            pos := components.PositionComponent.Get(entity)
            cx, cy := g.layout.HexToPixel(pos.Q, pos.R)
            vector.StrokeCircle(screen, float32(cx), float32(cy), hex.HexSize*0.9, 3,
                color.RGBA{255, 255, 0, 150}, false)
        }
    }

    // Draw UI
    g.drawUI(screen)

    // Draw victory/defeat screen
    if g.battle.State == states.Victory || g.battle.State == states.Defeat {
        g.drawGameOverScreen(screen)
    }
}

func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
    // Semi-transparent overlay
    overlay := ebiten.NewImage(screenWidth, screenHeight)
    overlay.Fill(color.RGBA{0, 0, 0, 180})
    screen.DrawImage(overlay, nil)

    // Message
    var message string
    var messageColor color.Color

    if g.battle.State == states.Victory {
        message = "VICTORY!"
        messageColor = color.RGBA{100, 255, 100, 255}
    } else {
        message = "DEFEAT"
        messageColor = color.RGBA{255, 100, 100, 255}
    }

    // Draw large message (we'd need a better text rendering solution for real styling)
    ebitenutil.DebugPrintAt(screen, message, screenWidth/2-50, screenHeight/2-50)
    ebitenutil.DebugPrintAt(screen, "Press R to restart", screenWidth/2-80, screenHeight/2)
}

func (g *Game) drawUI(screen *ebiten.Image) {
    // Instructions
    instructions := "Boss Battler - Demo Battle\n"
    if g.battle.State == states.WaitingForPlayerInput {
        instructions += "Click enemy to attack | W to wait\n"
    } else if g.battle.State == states.EnemyTurn {
        instructions += "Enemy's turn...\n"
    } else {
        instructions += fmt.Sprintf("%s\n", g.battle.State.String())
    }
    instructions += "ESC: quit"
    ebitenutil.DebugPrint(screen, instructions)

    // Turn order
    g.drawTurnOrderUI(screen)

    // Combat log
    g.drawCombatLog(screen)
}
```

### Step 6: Playtesting & Balance

Now run the game and play through multiple battles:

```bash
go run cmd/game/main.go
```

**Things to test**:
- [ ] Can you win by defeating the boss?
- [ ] Can you lose if all party members die?
- [ ] Are damage values reasonable?
- [ ] Is the boss too strong/weak?
- [ ] Do HP bars update correctly?
- [ ] Does combat log show all actions?
- [ ] Can you restart after game over?

**Balance adjustments** (tweak in entity creation files):
- Increase party HP if battles are too hard
- Decrease boss HP if battles are too long
- Adjust damage dice if fights feel too random
- Add more party members if it's too difficult

### Step 7: Final Polish

**Visual improvements**:
```go
// Draw actual hexagons instead of circles for entities
// (Use the hex corner drawing code from milestone 2)

// Add color flash when entity is hit
// Store "flash" timer in a component or map

// Shake screen on critical hit
// Offset all rendering by small random amount for 1-2 frames
```

**UI improvements**:
```go
// Show entity stats on hover
// Display weapon name and damage in tooltip
// Show action preview ("This will deal ~5-12 damage")
```

**Sound effects** (optional, requires sound library):
```bash
go get github.com/hajimehoshi/ebiten/v2/audio
```

Add simple beep sounds for hits/misses.

## Completion Checklist

When you've completed Milestone 7, you should have:

- [ ] Complete playable battle from start to finish
- [ ] Victory when all enemies defeated
- [ ] Defeat when all players defeated
- [ ] Restart functionality
- [ ] Smart boss AI (targets low HP)
- [ ] Visual feedback (HP bars, damage numbers)
- [ ] Improved UI
- [ ] Balanced gameplay
- [ ] No major bugs
- [ ] Fun to play!

## Congratulations!

You've built a complete boss-battler game from scratch! You now have:

1. **Solid foundation**: Ebiten + Donburi ECS + Hex Grid
2. **Complete systems**: Combat, turns, initiative, actions
3. **Player interaction**: Mouse input, target selection
4. **Game loop**: Start → Battle → Victory/Defeat → Restart
5. **Polish**: UI, feedback, balance

## Next Steps (Beyond the Guide)

Now that you have a working game, here are ideas for expansion:

### Short-term Enhancements
- [ ] Add more character classes
- [ ] Create different boss types
- [ ] Add special abilities (area attacks, healing, buffs)
- [ ] Implement movement system (click to move, then attack)
- [ ] Add ranged attack range limits
- [ ] Create multiple battle scenarios

### Medium-term Features
- [ ] Character progression (leveling, XP)
- [ ] Equipment system (weapons, armor, accessories)
- [ ] Status effects (poison, stun, buff, debuff)
- [ ] Party composition screen (choose your team)
- [ ] Multiple enemy types
- [ ] Procedural boss generation

### Long-term (Dungeon Crawler)
- [ ] Map system with rooms and corridors
- [ ] Random encounters
- [ ] Treasure and loot
- [ ] Saving/loading game state
- [ ] Multiple floors/dungeons
- [ ] Town system (shop, inn, quests)
- [ ] Story and dialogue

## Resources for Continued Learning

**Game Design**:
- D&D 5E Basic Rules (free): https://dnd.wizards.com/resources/systems-reference-document
- Red Blob Games: https://www.redblobgames.com/
- Game Programming Patterns: https://gameprogrammingpatterns.com/

**Ebiten Advanced Topics**:
- Shaders: https://ebitengine.org/en/documents/shader.html
- Audio: https://ebitengine.org/en/documents/audio.html
- Mobile: https://ebitengine.org/en/documents/mobile.html

**Go Game Dev Community**:
- Ebiten Discord: https://discord.gg/ebiten
- r/gamedev: https://reddit.com/r/gamedev

## Thank You!

You've completed all 7 milestones. I hope you learned a lot about Go, game development, ECS architecture, and had fun building this game!

Feel free to share your game, modify it, and make it your own. Good luck with your future game development projects!

---

*This guide was created to help you learn by doing. If you found it helpful, consider contributing your improvements back to help others learn too.*
