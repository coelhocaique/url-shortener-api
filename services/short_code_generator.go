package services

import (
	"crypto/rand"
	"encoding/hex"
)

// ShortCodeGenerator handles the generation of short codes
type ShortCodeGenerator struct{}

// NewShortCodeGenerator creates a new instance of ShortCodeGenerator
func NewShortCodeGenerator() *ShortCodeGenerator {
	return &ShortCodeGenerator{}
}

// Generate creates a random 5-character short code
func (g *ShortCodeGenerator) Generate() string {
	bytes := make([]byte, 3) // 3 bytes = 6 hex characters
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:5] // Take first 5 characters
}
