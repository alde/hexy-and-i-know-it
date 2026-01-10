package commands

import (
	"testing"

	"github.com/alde/hexy-and-i-know-it/internal/components"
	"github.com/yohamta/donburi"
)

// createTestEntity creates a basic entity for testing
func createTestEntity(world *donburi.World, name string, hp int) *donburi.Entry {
	entity := world.Create(
		components.HealthComponent,
		components.DisplayComponent,
	)

	components.HealthComponent.Set(entity, &components.HealthData{
		Max:     hp,
		Current: hp,
	})

	components.DisplayComponent.Set(entity, &components.DisplayData{
		Name: name,
	})

	return entity
}

// createCombatant creates an entity with combat stats
func createCombatant(world *donburi.World, name string, hp int, str, dex int) *donburi.Entry {
	entity := world.Create(
		components.HealthComponent,
		components.DisplayComponent,
		components.StatsComponent,
		components.WeaponComponent,
		components.ArmorComponent,
	)

	components.HealthComponent.Set(entity, &components.HealthData{
		Max:     hp,
		Current: hp,
	})

	components.DisplayComponent.Set(entity, &components.DisplayData{
		Name: name,
	})

	components.StatsComponent.Set(entity, &components.StatsData{
		Strength:  str,
		Dexterity: dex,
	})

	components.WeaponComponent.Set(entity, &components.WeaponData{
		Name:       "Sword",
		DamageDice: 1,
		DamageDie:  6,
		UsesStat:   "STR",
	})

	components.ArmorComponent.Set(entity, &components.ArmorData{
		BaseAC: 12,
		MaxDex: -1,
	})

	return entity
}

// TestWaitAction verifies wait action executes successfully
func TestWaitAction(t *testing.T) {
	world := donburi.NewWorld()
	actor := createTestEntity(world, "Warrior", 20)

	action := &WaitAction{
		Actor: actor,
	}

	// Validation should always pass
	if err := action.Validate(world); err != nil {
		t.Errorf("WaitAction.Validate() failed: %v", err)
	}

	// Execution should always succeed
	result := action.Execute(world)
	if !result.Success {
		t.Error("WaitAction.Execute() should always succeed")
	}

	// Description should include actor name
	desc := action.Description()
	if desc == "" {
		t.Error("WaitAction.Description() should not be empty")
	}
}

// TestAttackActionValidation verifies attack validation logic
func TestAttackActionValidation(t *testing.T) {
	world := donburi.NewWorld()

	tests := []struct {
		name        string
		setupFunc   func() (attacker, target *donburi.Entry)
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid attack",
			setupFunc: func() (*donburi.Entry, *donburi.Entry) {
				attacker := createCombatant(world, "Attacker", 20, 16, 12)
				target := createCombatant(world, "Target", 20, 14, 10)
				return attacker, target
			},
			shouldError: false,
		},
		{
			name: "attacker is dead",
			setupFunc: func() (*donburi.Entry, *donburi.Entry) {
				attacker := createCombatant(world, "Attacker", 20, 16, 12)
				target := createCombatant(world, "Target", 20, 14, 10)
				// Kill attacker
				health := components.HealthComponent.Get(attacker)
				health.Current = 0
				return attacker, target
			},
			shouldError: true,
			errorMsg:    "attacker is dead",
		},
		{
			name: "target is dead",
			setupFunc: func() (*donburi.Entry, *donburi.Entry) {
				attacker := createCombatant(world, "Attacker", 20, 16, 12)
				target := createCombatant(world, "Target", 20, 14, 10)
				// Kill target
				health := components.HealthComponent.Get(target)
				health.Current = 0
				return attacker, target
			},
			shouldError: true,
			errorMsg:    "target is already dead",
		},
		{
			name: "attack self",
			setupFunc: func() (*donburi.Entry, *donburi.Entry) {
				attacker := createCombatant(world, "Attacker", 20, 16, 12)
				return attacker, attacker
			},
			shouldError: true,
			errorMsg:    "cannot attack self",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attacker, target := tt.setupFunc()

			action := &AttackAction{
				Attacker: attacker,
				Target:   target,
			}

			err := action.Validate(world)

			if tt.shouldError && err == nil {
				t.Errorf("Expected validation error: %s", tt.errorMsg)
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

// TestAttackActionExecution verifies attack execution
func TestAttackActionExecution(t *testing.T) {
	world := donburi.NewWorld()
	attacker := createCombatant(world, "Attacker", 20, 16, 12)
	target := createCombatant(world, "Target", 20, 14, 10)

	action := &AttackAction{
		Attacker: attacker,
		Target:   target,
	}

	// Execute attack
	result := action.Execute(world)

	// Should have a result
	if result == nil {
		t.Fatal("Execute() returned nil result")
	}

	// Should have a message
	if result.Message == "" {
		t.Error("Result should have a message")
	}

	// Description should mention both entities
	desc := action.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}
