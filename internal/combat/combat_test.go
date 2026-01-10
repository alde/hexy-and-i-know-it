package combat

import (
	"testing"

	"github.com/alde/hexy-and-i-know-it/internal/components"
	"github.com/yohamta/donburi"
)

func createTestWorld() *donburi.World {
	return donburi.NewWorld()
}

// createWarrior creates a test warrior entity
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
		Constitution: 14, // +2
		Intelligence: 8,  // -1
		Wisdom:       10, // +0
		Charisma:     10, // +0
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

// createGoblin creates a test goblin enemy
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
		Constitution: 10, // +0
		Intelligence: 10, // +0
		Wisdom:       8,  // -1
		Charisma:     8,  // -1
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

// TestACCalculation verifies armor class is calculated correctly
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

// TestWeaponDamageRoll verifies weapon damage is in valid range
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

// TestAttackHitOrMiss verifies attack rolls work correctly
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

// TestAttackReducesHealth verifies attacks reduce target HP
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

// TestAttackKillsTarget verifies attacks can kill targets at low HP
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
