package hex

import (
	"math"
	"testing"
)

// TestHexToPixel verifies hex coordinates convert to pixel coordinates
func TestHexToPixel(t *testing.T) {
	layout := NewLayout()

	tests := []struct {
		name string
		q, r int64
	}{
		{"origin", 0, 0},
		{"positive q", 1, 0},
		{"positive r", 0, 1},
		{"negative q", -1, 0},
		{"negative r", 0, -1},
		{"diagonal", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := layout.HexToPixel(tt.q, tt.r)

			// Basic sanity checks
			if math.IsNaN(x) || math.IsNaN(y) {
				t.Errorf("HexToPixel(%d, %d) returned NaN", tt.q, tt.r)
			}
			if math.IsInf(x, 0) || math.IsInf(y, 0) {
				t.Errorf("HexToPixel(%d, %d) returned Inf", tt.q, tt.r)
			}
		})
	}
}

// TestPixelToHex verifies pixel coordinates convert back to hex coordinates
func TestPixelToHex(t *testing.T) {
	layout := NewLayout()

	tests := []struct {
		name         string
		x, y         float64
		wantQ, wantR int64
	}{
		{"center of screen", OffsetX, OffsetY, 0, 0},
		// Add more tests after you verify the first one works
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r := layout.PixelToHex(tt.x, tt.y)

			if q != tt.wantQ || r != tt.wantR {
				t.Errorf("PixelToHex(%f, %f) = (%d, %d), want (%d, %d)",
					tt.x, tt.y, q, r, tt.wantQ, tt.wantR)
			}
		})
	}
}

// TestRoundTrip verifies hex -> pixel -> hex returns the original hex
func TestRoundTrip(t *testing.T) {
	layout := NewLayout()

	tests := []struct {
		q, r int64
	}{
		{0, 0},
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
		{3, 3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Hex to pixel
			x, y := layout.HexToPixel(tt.q, tt.r)

			// Pixel back to hex
			q, r := layout.PixelToHex(x, y)

			if q != tt.q || r != tt.r {
				t.Errorf("Round trip failed: (%d,%d) -> (%f,%f) -> (%d,%d)",
					tt.q, tt.r, x, y, q, r)
			}
		})
	}
}
