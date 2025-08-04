package random

import (
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		chars  Chars
		want   int
	}{
		{"empty length", 0, AlphaLarge | AlphaSmall | ZeroNine, 0},
		{"no charset", 10, 0, 0},
		{"only large", 10, AlphaLarge, 10},
		{"only small", 10, AlphaSmall, 10},
		{"only numbers", 10, ZeroNine, 10},
		{"large and small", 10, AlphaLarge | AlphaSmall, 10},
		{"all chars", 20, AlphaLarge | AlphaSmall | ZeroNine, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.length, tt.chars)
			if len(got) != tt.want {
				t.Errorf("String() length = %v, want %v", len(got), tt.want)
			}

			if tt.length > 0 && tt.chars > 0 {
				// Verify character set presence
				for _, c := range got {
					valid := false
					if tt.chars&AlphaLarge != 0 && c >= 'A' && c <= 'Z' {
						valid = true
					}
					if tt.chars&AlphaSmall != 0 && c >= 'a' && c <= 'z' {
						valid = true
					}
					if tt.chars&ZeroNine != 0 && c >= '0' && c <= '9' {
						valid = true
					}
					if !valid {
						t.Errorf("String() contains invalid character: %c", c)
					}
				}
			}
		})
	}
}
