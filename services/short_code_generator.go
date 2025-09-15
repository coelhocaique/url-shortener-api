package services

import (
	"url-shortener-api/models"
)

// ShortCodeGenerator handles the generation of short codes using distributed counter
type ShortCodeGenerator struct {
	counter models.CounterService
	encoder *Base62Encoder
}

// NewShortCodeGenerator creates a new instance of ShortCodeGenerator
func NewShortCodeGenerator(counter models.CounterService) *ShortCodeGenerator {
	return &ShortCodeGenerator{
		counter: counter,
		encoder: NewBase62Encoder(),
	}
}

// Generate creates a short code using the distributed counter and base62 encoding
func (g *ShortCodeGenerator) Generate() (string, error) {
	counter, err := g.counter.GetNextCounter()
	if err != nil {
		return "", err
	}

	return g.encoder.Encode(counter), nil
}
