package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "     hello      world		",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    " goto this-city-OVER-here	@1.1",
			expected: []string{"goto", "this-city-over-here", "@1.1"},
		},
	}

	for i, c := range cases {
		t.Logf("Test case #%d", i+1)
		actual := cleanInput(c.input)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("Expected %v, actual %v", c.expected, actual)
		}
	}
}
