package main

import "testing"

func eval(input string) (int, error) {
	parser := &Parser{input: input}
	return parser.parseExpression()
}

func TestCalc(t *testing.T) {
	tests := []struct {
		expr     string
		expected int
		wantErr  bool
	}{
		{"1+2", 3, false},
		{"6-2", 4, false},
		{"3*2", 6, false},
		{"9/3", 3, false},
		{"(1+2)-3", 0, false},
		{"(1+2)*4", 12, false},
		{"2*(3+4)", 14, false},
		{"(2+6)*(3+1)", 32, false},
		{"10-2*3", 4, false},
		{"(12-2)*3", 30, false},
		{"15/3+2", 7, false},
		{"24/(1+3)", 6, false},

		{"-5", -5, false},
		{"-(2+3)", -5, false},
		{"-3*2", -6, false},
		{"2+-3", -1, false},

		{"5/0", 0, true},
		{"(2+3", 0, true},
		{"6-*2", 0, true},
		{"xyz", 0, true},
		{"6/(3-3)", 0, true},
	}

	for _, test := range tests {
		got, err := eval(test.expr)
		if test.wantErr {
			if err == nil {
				t.Errorf("Для '%s' ожмдалась ошибка, но не возникла", test.expr)
			}
			continue
		}
		if err != nil {
			t.Errorf("Для '%s' не ожидалась ошибка, но возникла: %v", test.expr, err)
			continue
		}
		if got != test.expected {
			t.Errorf("Для '%s' ожидалось %d, получили %d", test.expr, test.expected, got)
		}
	}
}
