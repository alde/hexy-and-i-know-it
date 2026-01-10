package components

import (
	"testing"
)

// TestStatModifier verifies D&D stat modifier calculation
// Formula: (stat - 10) / 2, rounded down
func TestStatModifier(t *testing.T) {
	tests := []struct {
		name     string
		stat     int
		expected int
	}{
		{"very low (3)", 3, -4},
		{"low (8)", 8, -1},
		{"average (10)", 10, 0},
		{"good (14)", 14, 2},
		{"great (18)", 18, 4},
		{"legendary (20)", 20, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := StatsData{Strength: tt.stat}
			got := stats.Modifier(tt.stat)

			if got != tt.expected {
				t.Errorf("Modifier(%d) = %d, want %d", tt.stat, got, tt.expected)
			}
		})
	}
}

// TestHealthDamage verifies damage reduces HP correctly
func TestHealthDamage(t *testing.T) {
	health := HealthData{Current: 50, Max: 100}

	health.Damage(20)
	if health.Current != 30 {
		t.Errorf("After 20 damage: got %d HP, want 30", health.Current)
	}

	health.Damage(100)
	if health.Current != 0 {
		t.Errorf("After overkill damage: got %d HP, want 0", health.Current)
	}

	if !health.IsDead() {
		t.Error("Entity should be dead at 0 HP")
	}
}

// TestHealthHeal verifies healing increases HP correctly
func TestHealthHeal(t *testing.T) {
	health := HealthData{Current: 30, Max: 100}

	health.Heal(20)
	if health.Current != 50 {
		t.Errorf("After healing 20: got %d HP, want 50", health.Current)
	}

	health.Heal(100)
	if health.Current != 100 {
		t.Errorf("After overheal: got %d HP, want 100 (max)", health.Current)
	}

	if !health.IsFullHealth() {
		t.Error("Entity should be at full health")
	}
}

// TestSizeNumHexes verifies hex occupation calculation
// Formula for hex rings: 1 + 3*n*(n+1) where n is radius
func TestSizeNumHexes(t *testing.T) {
	tests := []struct {
		radius   int
		expected int
	}{
		{0, 1},   // Single hex
		{1, 7},   // Center + 6 surrounding
		{2, 19},  // Center + 6 + 12
		{3, 37},  // Center + 6 + 12 + 18
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			size := SizeData{Radius: tt.radius}
			got := size.NumHexes()

			if got != tt.expected {
				t.Errorf("Radius %d: got %d hexes, want %d", tt.radius, got, tt.expected)
			}
		})
	}
}
