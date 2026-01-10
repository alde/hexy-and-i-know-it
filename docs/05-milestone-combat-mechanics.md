# Milestone 5: Combat Mechanics

## Goal

Implement D&D-style combat mechanics: attack rolls, damage calculation, armor class, and effects.

## Learning Objectives

- Understand D&D 5E combat mechanics
- Implement attack rolls (d20 + modifiers vs AC)
- Calculate damage (dice rolls + modifiers)
- Add armor class (AC) component
- Create combat system for resolving attacks
- Handle critical hits and misses

## Prerequisites

- Completed Milestone 4
- Basic understanding of D&D mechanics

## Tasks

### Combat Components
- [ ] Create `internal/components/armor.go` for AC (Armor Class)
- [ ] Create `internal/components/weapon.go` for damage dice
- [ ] Update entities to have armor and weapons
- [ ] Add attack bonus calculation

### Combat System
- [ ] Create `internal/combat/attack.go` for attack resolution
- [ ] Implement attack roll (d20 + attack bonus vs AC)
- [ ] Implement damage roll (dice + ability modifier)
- [ ] Handle critical hits (nat 20)
- [ ] Handle critical misses (nat 1)
- [ ] Create tests in `internal/combat/combat_test.go`

### Damage Application
- [ ] Create function to apply damage to entities
- [ ] Update health component
- [ ] Handle entity death
- [ ] Create combat log component

### Integration
- [ ] Add simple attack action (press A to attack)
- [ ] Show attack results in UI
- [ ] Display damage numbers
- [ ] Remove dead entities from battle

### Testing
- [ ] Run tests: `go test ./internal/combat`
- [ ] Attack rolls work correctly
- [ ] Damage calculated properly
- [ ] Critical hits deal double damage
- [ ] AC blocks some attacks
- [ ] Entities die at 0 HP

## D&D 5E Combat Primer

### Attack Roll

**To Hit Roll**:
```
d20 + Attack Bonus >= Target AC
```

- **d20**: 20-sided die (1-20)
- **Attack Bonus**: Ability modifier + proficiency bonus
- **Target AC**: Armor Class (10 + modifiers)

**Result**:
- Natural 20: Critical hit (double damage dice)
- Natural 1: Critical miss (always fails)
- Roll >= AC: Hit
- Roll < AC: Miss

### Damage Roll

When you hit:
```
Weapon Damage Dice + Ability Modifier
```

Examples:
- Longsword: 1d8 + STR modifier
- Dagger: 1d4 + DEX modifier (finesse)
- Fireball: 8d6 (no modifier for area spells)

**Critical Hit**:
- Double the dice (not the modifier)
- 1d8+3 becomes 2d8+3

### Armor Class (AC)

Base AC calculation:
```
AC = 10 + DEX modifier + armor bonus
```

Examples:
- Unarmored: 10 + DEX (cloth)
- Light armor: 11 + DEX (leather)
- Medium armor: 13 + min(DEX, 2) (chain shirt)
- Heavy armor: 16 + 0 (plate, no DEX)

### Attack Bonus

Melee attack:
```
d20 + STR modifier + proficiency
```

Ranged attack:
```
d20 + DEX modifier + proficiency
```

For level 1, proficiency = +2

## Step-by-Step Implementation

### Step 1: Create Armor Component

Create `internal/components/armor.go`:

```go
package components

import "github.com/yohamta/donburi"

// ArmorData stores armor class information
type ArmorData struct {
	BaseAC int  // Base armor value
	MaxDex int  // Max DEX bonus allowed (-1 = unlimited)
}

// CalculateAC computes final AC with DEX modifier
func (a *ArmorData) CalculateAC(dexMod int) int {
	ac := a.BaseAC

	if a.MaxDex == -1 {
		// Unlimited DEX bonus (light/no armor)
		ac += dexMod
	} else if a.MaxDex > 0 {
		// Limited DEX bonus (medium armor)
		if dexMod > a.MaxDex {
			ac += a.MaxDex
		} else {
			ac += dexMod
		}
	}
	// else MaxDex == 0: heavy armor, no DEX bonus

	return ac
}

var ArmorComponent = donburi.NewComponentType[ArmorData]()
```

Create `internal/components/weapon.go`:

```go
package components

import (
	"math/rand"

	"github.com/yohamta/donburi"
)

// WeaponData stores weapon information
type WeaponData struct {
	Name       string
	DamageDice int    // Number of dice (e.g., 1 for 1d8)
	DamageDie  int    // Die size (e.g., 8 for 1d8)
	UsesStat   string // "STR" or "DEX" for attack and damage bonus
}

// RollDamage rolls weapon damage dice
func (w *WeaponData) RollDamage() int {
	total := 0
	for i := 0; i < w.DamageDice; i++ {
		total += rand.Intn(w.DamageDie) + 1
	}
	return total
}

var WeaponComponent = donburi.NewComponentType[WeaponData]()
```

### Step 2: Create Combat System

Create `internal/combat/attack.go`:

```go
package combat

import (
	"fmt"
	"math/rand"

	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
)

// AttackResult represents the outcome of an attack
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
}

// String formats the attack result for display
func (r *AttackResult) String() string {
	if !r.Hit {
		return fmt.Sprintf("%s attacks %s... MISS! (rolled %d vs AC %d)",
			r.AttackerName, r.TargetName, r.TotalAttack, r.TargetAC)
	}

	critText := ""
	if r.Critical {
		critText = " CRITICAL HIT!"
	}

	killText := ""
	if r.Killed {
		killText = " TARGET SLAIN!"
	}

	return fmt.Sprintf("%s attacks %s... HIT!%s (rolled %d vs AC %d) for %d damage%s",
		r.AttackerName, r.TargetName, critText, r.TotalAttack, r.TargetAC, r.Damage, killText)
}

const (
	ProficiencyBonus = 2 // Level 1 proficiency
)

// PerformAttack executes an attack from attacker to target
func PerformAttack(attacker, target *donburi.Entry) *AttackResult {
	result := &AttackResult{}

	// Get attacker info
	attackerDisplay := components.DisplayComponent.Get(attacker)
	attackerStats := components.StatsComponent.Get(attacker)
	attackerWeapon := components.WeaponComponent.Get(attacker)

	// Get target info
	targetDisplay := components.DisplayComponent.Get(target)
	targetStats := components.StatsComponent.Get(target)
	targetArmor := components.ArmorComponent.Get(target)
	targetHealth := components.HealthComponent.Get(target)

	result.AttackerName = attackerDisplay.Name
	result.TargetName = targetDisplay.Name

	// Calculate target AC
	result.TargetAC = targetArmor.CalculateAC(targetStats.DexMod())

	// Determine attack bonus
	var attackBonus int
	if attackerWeapon.UsesStat == "STR" {
		attackBonus = attackerStats.StrMod()
	} else if attackerWeapon.UsesStat == "DEX" {
		attackBonus = attackerStats.DexMod()
	}
	attackBonus += ProficiencyBonus

	// Roll attack (d20)
	result.AttackRoll = rand.Intn(20) + 1
	result.TotalAttack = result.AttackRoll + attackBonus

	// Check for critical hit/miss
	if result.AttackRoll == 20 {
		result.Critical = true
		result.Hit = true
	} else if result.AttackRoll == 1 {
		result.Hit = false
		return result
	} else {
		result.Hit = result.TotalAttack >= result.TargetAC
	}

	// If hit, roll damage
	if result.Hit {
		baseDamage := attackerWeapon.RollDamage()

		// Critical hit doubles dice (not modifier)
		if result.Critical {
			baseDamage += attackerWeapon.RollDamage()
		}

		// Add ability modifier
		var damageBonus int
		if attackerWeapon.UsesStat == "STR" {
			damageBonus = attackerStats.StrMod()
		} else if attackerWeapon.UsesStat == "DEX" {
			damageBonus = attackerStats.DexMod()
		}

		result.Damage = baseDamage + damageBonus
		if result.Damage < 1 {
			result.Damage = 1 // Minimum 1 damage on hit
		}

		// Apply damage
		targetHealth.Damage(result.Damage)

		// Check if killed
		if targetHealth.IsDead() {
			result.Killed = true
		}
	}

	return result
}
```

### Step 3: Create Combat Tests

Create `internal/combat/combat_test.go`:

```go
package combat

import (
	"testing"

	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
)

func createTestWorld() *donburi.World {
	return donburi.NewWorld()
}

func createWarrior(world *donburi.World) *donburi.Entry {
	entity := world.Create(
		components.StatsComponent,
		components.HealthComponent,
		components.WeaponComponent,
		components.ArmorComponent,
		components.DisplayComponent,
	)

	components.StatsComponent.Set(entity, &components.StatsData{
		Strength:     16, // +3
		Dexterity:    12, // +1
		Constitution: 14,
		Intelligence: 8,
		Wisdom:       10,
		Charisma:     10,
	})

	components.HealthComponent.Set(entity, &components.HealthData{
		Max:     20,
		Current: 20,
	})

	components.WeaponComponent.Set(entity, &components.WeaponData{
		Name:       "Longsword",
		DamageDice: 1,
		DamageDie:  8,
		UsesStat:   "STR",
	})

	components.ArmorComponent.Set(entity, &components.ArmorData{
		BaseAC: 15, // Chain mail
		MaxDex: 0,  // Heavy armor
	})

	components.DisplayComponent.Set(entity, &components.DisplayData{
		Name: "Test Warrior",
	})

	return entity
}

func createGoblin(world *donburi.World) *donburi.Entry {
	entity := world.Create(
		components.StatsComponent,
		components.HealthComponent,
		components.WeaponComponent,
		components.ArmorComponent,
		components.DisplayComponent,
	)

	components.StatsComponent.Set(entity, &components.StatsData{
		Strength:     8,  // -1
		Dexterity:    14, // +2
		Constitution: 10,
		Intelligence: 10,
		Wisdom:       8,
		Charisma:     8,
	})

	components.HealthComponent.Set(entity, &components.HealthData{
		Max:     7,
		Current: 7,
	})

	components.WeaponComponent.Set(entity, &components.WeaponData{
		Name:       "Scimitar",
		DamageDice: 1,
		DamageDie:  6,
		UsesStat:   "DEX",
	})

	components.ArmorComponent.Set(entity, &components.ArmorData{
		BaseAC: 12, // Leather
		MaxDex: -1, // Light armor, unlimited DEX
	})

	components.DisplayComponent.Set(entity, &components.DisplayData{
		Name: "Test Goblin",
	})

	return entity
}

func TestACCalculation(t *testing.T) {
	tests := []struct {
		name   string
		armor  components.ArmorData
		dexMod int
		wantAC int
	}{
		{"unarmored high DEX", components.ArmorData{BaseAC: 10, MaxDex: -1}, 3, 13},
		{"light armor", components.ArmorData{BaseAC: 11, MaxDex: -1}, 2, 13},
		{"medium armor capped", components.ArmorData{BaseAC: 13, MaxDex: 2}, 4, 15},
		{"medium armor not capped", components.ArmorData{BaseAC: 13, MaxDex: 2}, 1, 14},
		{"heavy armor", components.ArmorData{BaseAC: 16, MaxDex: 0}, 3, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAC := tt.armor.CalculateAC(tt.dexMod)
			if gotAC != tt.wantAC {
				t.Errorf("CalculateAC() = %d, want %d", gotAC, tt.wantAC)
			}
		})
	}
}

func TestWeaponDamageRoll(t *testing.T) {
	weapon := components.WeaponData{
		Name:       "Greatsword",
		DamageDice: 2,
		DamageDie:  6,
		UsesStat:   "STR",
	}

	// Roll multiple times to verify range
	for i := 0; i < 100; i++ {
		damage := weapon.RollDamage()
		if damage < 2 || damage > 12 {
			t.Errorf("RollDamage() = %d, want range [2, 12]", damage)
		}
	}
}

func TestAttackHitOrMiss(t *testing.T) {
	world := createTestWorld()
	warrior := createWarrior(world)
	goblin := createGoblin(world)

	// Run multiple attacks to test randomness
	hits := 0
	misses := 0

	for i := 0; i < 100; i++ {
		// Reset goblin health
		goblinHealth := components.HealthComponent.Get(goblin)
		goblinHealth.Current = goblinHealth.Max

		result := PerformAttack(warrior, goblin)

		if result.Hit {
			hits++
			if result.Damage < 1 {
				t.Error("Hit should deal at least 1 damage")
			}
		} else {
			misses++
		}
	}

	// Warrior has +5 attack (+3 STR, +2 prof) vs AC 14 (12 base + 2 DEX)
	// Should hit on 9+ (60% of the time)
	if hits < 40 || hits > 80 {
		t.Logf("Warning: Hit rate seems off. Hits: %d, Misses: %d", hits, misses)
		// Not a hard failure, randomness can vary
	}
}

func TestAttackReducesHealth(t *testing.T) {
	world := createTestWorld()
	warrior := createWarrior(world)
	goblin := createGoblin(world)

	goblinHealth := components.HealthComponent.Get(goblin)
	initialHP := goblinHealth.Current

	result := PerformAttack(warrior, goblin)

	if result.Hit {
		if goblinHealth.Current >= initialHP {
			t.Error("Health should decrease after being hit")
		}

		expectedHP := initialHP - result.Damage
		if expectedHP < 0 {
			expectedHP = 0
		}

		if goblinHealth.Current != expectedHP {
			t.Errorf("Health = %d, expected %d", goblinHealth.Current, expectedHP)
		}
	}
}

func TestAttackKillsTarget(t *testing.T) {
	world := createTestWorld()
	warrior := createWarrior(world)
	goblin := createGoblin(world)

	// Set goblin to very low HP
	goblinHealth := components.HealthComponent.Get(goblin)
	goblinHealth.Current = 1

	// Keep attacking until we get a hit
	for i := 0; i < 100; i++ {
		result := PerformAttack(warrior, goblin)
		if result.Hit {
			if !result.Killed {
				t.Error("Goblin with 1 HP should be killed by any hit")
			}
			if goblinHealth.Current != 0 {
				t.Errorf("Dead entity should have 0 HP, got %d", goblinHealth.Current)
			}
			return
		}
		// Reset for next attempt
		goblinHealth.Current = 1
	}

	t.Error("Never landed a hit in 100 attempts (very unlikely)")
}
```

Run tests:
```bash
go test ./internal/combat
```

### Step 4: Update Entities with Combat Components

Update `internal/entities/character.go`:

```go
// In SpawnCharacter, add to entity.Create:
entity := world.Create(
    components.PositionComponent,
    components.StatsComponent,
    components.HealthComponent,
    components.DisplayComponent,
    components.SizeComponent,
    components.InitiativeComponent,
    components.ArmorComponent,  // ADD
    components.WeaponComponent, // ADD
)

// After setting display, add weapon and armor based on class:

switch class {
case Warrior:
    // ... existing stats and display code ...

    components.WeaponComponent.Set(entity, &components.WeaponData{
        Name:       "Longsword",
        DamageDice: 1,
        DamageDie:  8,
        UsesStat:   "STR",
    })

    components.ArmorComponent.Set(entity, &components.ArmorData{
        BaseAC: 15, // Chain mail
        MaxDex: 0,
    })

case Rogue:
    components.WeaponComponent.Set(entity, &components.WeaponData{
        Name:       "Rapier",
        DamageDice: 1,
        DamageDie:  8,
        UsesStat:   "DEX",
    })

    components.ArmorComponent.Set(entity, &components.ArmorData{
        BaseAC: 12, // Leather
        MaxDex: -1,
    })

case Mage:
    components.WeaponComponent.Set(entity, &components.WeaponData{
        Name:       "Dagger",
        DamageDice: 1,
        DamageDie:  4,
        UsesStat:   "DEX",
    })

    components.ArmorComponent.Set(entity, &components.ArmorData{
        BaseAC: 10, // Robes (no armor)
        MaxDex: -1,
    })

case Cleric:
    components.WeaponComponent.Set(entity, &components.WeaponData{
        Name:       "Mace",
        DamageDice: 1,
        DamageDie:  6,
        UsesStat:   "STR",
    })

    components.ArmorComponent.Set(entity, &components.ArmorData{
        BaseAC: 14, // Scale mail
        MaxDex: 2,  // Medium armor
    })
}
```

Update `internal/entities/boss.go` similarly:

```go
// In SpawnBoss, add components and set weapon/armor:
components.WeaponComponent.Set(entity, &components.WeaponData{
    Name:       "Claws",
    DamageDice: 2,
    DamageDie:  6,
    UsesStat:   "STR",
})

components.ArmorComponent.Set(entity, &components.ArmorData{
    BaseAC: 13, // Natural armor
    MaxDex: 2,
})
```

### Step 5: Test Combat in Game

Update `cmd/game/main.go` to add simple attack test:

```go
// In Update method, add test attack key:
if inpututil.IsKeyJustPressed(ebiten.KeyA) {
    // Get first two entities and make one attack the other (testing)
    query := donburi.NewQuery(
        filter.Contains(components.HealthComponent),
    )

    var entities []*donburi.Entry
    query.Each(g.world, func(entry *donburi.Entry) {
        entities = append(entities, entry)
    })

    if len(entities) >= 2 {
        result := combat.PerformAttack(entities[0], entities[1])
        log.Println(result.String())
    }
}

// Update UI text:
ebitenutil.DebugPrint(screen, "Milestone 5: Combat\nSPACE: next turn\nA: test attack\nESC: quit")
```

Don't forget imports:
```go
import (
	"log"
	"github.com/yourusername/boss-battler/internal/combat"
	"github.com/yohamta/donburi/filter"
)
```

Run the game and press A to see attack results in the console.

## Key Concepts

### D&D Attack Flow

1. Declare attack
2. Roll d20 + attack bonus
3. Compare to target AC
4. If hit, roll damage dice + modifier
5. Apply damage to target
6. Check if target dies

### Critical Hits

- Natural 20 on d20 always hits
- Double damage **dice** (not modifiers)
- Example: 1d8+3 crit = 2d8+3

### Minimum Damage

In this implementation, successful hits always deal at least 1 damage, even if modifiers are negative.

### Proficiency Bonus

Represents training/skill. Increases with level:
- Level 1-4: +2
- Level 5-8: +3
- Level 9-12: +4
- etc.

## Common Issues

### Damage always 0
- Check weapon damage die size
- Verify damage modifier is being added
- Ensure RollDamage is being called

### Too many hits/misses
- Check AC calculation
- Verify attack bonus includes proficiency
- Check d20 roll is 1-20 (not 0-19)

### AC seems wrong
- Verify DEX modifier calculation
- Check MaxDex limits for medium/heavy armor
- Print intermediate values to debug

### Critical hits not working
- Check for `== 20` (natural 20, not total)
- Verify damage dice are being doubled (not modifier)

## Next Steps

Milestone 5 complete! You now have:
- Full D&D-style combat mechanics
- Attack rolls vs AC
- Damage calculation with dice
- Critical hits
- Entity death

In [Milestone 6](06-milestone-input-actions.md), you'll add player input to select targets and execute attacks manually!

## Extra Challenges (Optional)

- [ ] Add advantage/disadvantage (roll 2d20, take higher/lower)
- [ ] Implement different damage types (slashing, piercing, bludgeoning)
- [ ] Add resistances/vulnerabilities
- [ ] Create spell attacks (INT/WIS + proficiency)
- [ ] Add saving throws
- [ ] Implement status effects (poisoned, stunned, etc.)
- [ ] Add healing spells/potions
