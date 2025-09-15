package services_test

import (
	"strings"
	"testing"
	"url-shortener-api/services"
	"url-shortener-api/tests/testutils"
)

func TestShortCodeGenerator_Generate(t *testing.T) {
	mockCounter := testutils.NewMockCounterService()
	generator := services.NewShortCodeGenerator(mockCounter)

	tests := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "generates base62 encoded code",
			expectedLength: 1, // First counter value (1) encoded in base62 is "1"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.Generate()
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}

			if len(result) != tt.expectedLength {
				t.Errorf("Generate() length = %d, want %d", len(result), tt.expectedLength)
			}

			// Check if result contains only base62 characters
			for _, char := range result {
				if !strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", char) {
					t.Errorf("Generate() contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestShortCodeGenerator_GenerateConsistentLength(t *testing.T) {
	mockCounter := testutils.NewMockCounterService()
	generator := services.NewShortCodeGenerator(mockCounter)
	lengths := make(map[int]int)

	for i := 0; i < 10; i++ {
		code, err := generator.Generate()
		if err != nil {
			t.Errorf("Generate() error = %v", err)
			return
		}
		lengths[len(code)]++
	}

	// All codes should be the same length (base62 encoding of sequential numbers)
	if len(lengths) != 1 {
		t.Errorf("Generate() produced codes with different lengths: %v", lengths)
	}

	// Check if the length is 1 (for first 10 counter values)
	for length := range lengths {
		if length != 1 {
			t.Errorf("Generate() produced code with length %d, want 1", length)
		}
	}
}
