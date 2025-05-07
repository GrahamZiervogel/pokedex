package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Basic case with leading/trailing spaces",
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Mixed case and multiple words",
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			name:     "Single word with leading/trailing spaces",
			input:    "  single  ",
			expected: []string{"single"},
		},
		{
			name:     "Single uppercase word",
			input:    "UPPERCASE",
			expected: []string{"uppercase"},
		},
		{
			name:     "Input with only whitespace characters",
			input:    "   \t  \n   ",
			expected: []string{},
		},
		{
			name:     "Empty string input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Mixed whitespace types between words",
			input:    "  one \t two  \n three  ",
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "No leading/trailing spaces, already lowercase",
			input:    "word",
			expected: []string{"word"},
		},
		{
			name:     "Multiple spaces between words",
			input:    "alpha   beta  gamma",
			expected: []string{"alpha", "beta", "gamma"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected %d words, got %d", len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Expected word %d to be '%s', got '%s'", i, expectedWord, word)
			}
		}
	}
}
