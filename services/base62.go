package services

import (
	"fmt"
	"math"
	"strings"
)

const (
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base62Length   = len(base62Alphabet)
)

// Base62Encoder handles base62 encoding operations
type Base62Encoder struct{}

// NewBase62Encoder creates a new instance of Base62Encoder
func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{}
}

// Encode converts a number to base62 string
func (e *Base62Encoder) Encode(num int64) string {
	if num == 0 {
		return "0"
	}

	var result strings.Builder
	for num > 0 {
		result.WriteByte(base62Alphabet[num%int64(base62Length)])
		num /= int64(base62Length)
	}

	// Reverse the string
	encoded := result.String()
	reversed := make([]byte, len(encoded))
	for i, j := 0, len(encoded)-1; i < len(encoded); i, j = i+1, j-1 {
		reversed[i] = encoded[j]
	}

	return string(reversed)
}

// Decode converts a base62 string to number
func (e *Base62Encoder) Decode(str string) (int64, error) {
	var result int64
	base := int64(base62Length)

	for i, char := range str {
		index := strings.IndexRune(base62Alphabet, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid character '%c' in base62 string", char)
		}
		result += int64(index) * int64(math.Pow(float64(base), float64(len(str)-1-i)))
	}

	return result, nil
}
