package services_test

import (
	"strings"
	"testing"
	"url-shortener-api/services"
)

func TestShortCodeGenerator_Generate(t *testing.T) {
	generator := services.NewShortCodeGenerator()

	tests := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "generates 5 character code",
			expectedLength: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.Generate()

			if len(result) != tt.expectedLength {
				t.Errorf("Generate() length = %d, want %d", len(result), tt.expectedLength)
			}

			// Check if result contains only hex characters
			for _, char := range result {
				if !strings.ContainsRune("0123456789abcdef", char) {
					t.Errorf("Generate() contains invalid character: %c", char)
				}
			}
		})
	}
}


func TestShortCodeGenerator_GenerateConsistentLength(t *testing.T) {
	generator := services.NewShortCodeGenerator()
	lengths := make(map[int]int)

	for i := 0; i < 100; i++ {
		code := generator.Generate()
		lengths[len(code)]++
	}

	// All codes should be the same length
	if len(lengths) != 1 {
		t.Errorf("Generate() produced codes with different lengths: %v", lengths)
	}

	// Check if the length is 5
	for length := range lengths {
		if length != 5 {
			t.Errorf("Generate() produced code with length %d, want 5", length)
		}
	}
}
