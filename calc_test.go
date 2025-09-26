package main

import "testing"

func eval(input string) int {
	parser := &Parser{input: input}
	return parser.parseExpression()
}

func TestCalc(t *testing.T) {
	tests := []struct {
		expr     string
		expected int
	}{
		{"1+2", 3},
		{"6-2", 4},
		{"3*2", 6},
		{"9/3", 3},
		{"(1+2)-3", 0},
		{"(1+2)*4", 12},
		{"2*(3+4)", 14},
		{"(2+6)*(3+1)", 32},
		{"10-2*3", 4},
		{"(12-2)*3", 30},
		{"15/3+2", 7},
		{"24/(1+3)", 6},
	}

	for _, test := range tests {
		result := eval(test.expr)
		if result != test.expected {
			t.Errorf("For '%s', expected %d, got %d", test.expr, test.expected, result)
		}
	}
}
