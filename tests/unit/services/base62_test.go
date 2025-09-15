package services_test

import (
	"testing"
	"url-shortener-api/services"
)

func TestBase62Encoder_Encode(t *testing.T) {
	encoder := services.NewBase62Encoder()

	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "encode 0",
			input:    0,
			expected: "0",
		},
		{
			name:     "encode 1",
			input:    1,
			expected: "1",
		},
		{
			name:     "encode 10",
			input:    10,
			expected: "A",
		},
		{
			name:     "encode 61",
			input:    61,
			expected: "z",
		},
		{
			name:     "encode 62",
			input:    62,
			expected: "10",
		},
		{
			name:     "encode 3844",
			input:    3844,
			expected: "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoder.Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Encode(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBase62Encoder_Decode(t *testing.T) {
	encoder := services.NewBase62Encoder()

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{
			name:     "decode 0",
			input:    "0",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "decode 1",
			input:    "1",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "decode A",
			input:    "A",
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "decode z",
			input:    "z",
			expected: 61,
			wantErr:  false,
		},
		{
			name:     "decode 10",
			input:    "10",
			expected: 62,
			wantErr:  false,
		},
		{
			name:     "decode 100",
			input:    "100",
			expected: 3844,
			wantErr:  false,
		},
		{
			name:     "decode invalid character",
			input:    "!",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Decode(%s) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBase62Encoder_RoundTrip(t *testing.T) {
	encoder := services.NewBase62Encoder()

	// Test round-trip encoding/decoding
	testValues := []int64{0, 1, 10, 61, 62, 100, 1000, 10000, 100000}

	for _, value := range testValues {
		encoded := encoder.Encode(value)
		decoded, err := encoder.Decode(encoded)
		if err != nil {
			t.Errorf("Decode failed for %d: %v", value, err)
			continue
		}
		if decoded != value {
			t.Errorf("Round-trip failed for %d: encoded to %s, decoded to %d", value, encoded, decoded)
		}
	}
}
