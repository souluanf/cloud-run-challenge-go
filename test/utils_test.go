package test

import (
	"cloud-run-challenge-go/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Olá, mundo!", "Ola, mundo!"},
		{"Coração", "Coracao"},
		{"ÀÉÌÕÙ", "AEIOU"},
		{"áéíóúç", "aeiouc"},
	}

	for _, tt := range tests {
		result := internal.RemoveAccents(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}
