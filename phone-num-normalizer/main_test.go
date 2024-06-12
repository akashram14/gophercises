package main

import "testing"

func TestNormalize(t *testing.T) {

	testCases := []struct {
		input    string
		expected string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.expected {
				t.Errorf("got %s: want: %s", actual, tc.expected)
			}
		})
	}

}
