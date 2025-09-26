package uniq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     Options
		expected string
	}{
		{
			name:     "basic",
			input:    "a\na\na\nb\n",
			opts:     Options{},
			expected: "a\nb",
		},
		{
			name:     "count",
			input:    "a\na\nb\n",
			opts:     Options{Count: true},
			expected: "2 a\n1 b",
		},
		{
			name:     "only dups",
			input:    "apple\napple\nbanana\n",
			opts:     Options{OnlyDup: true},
			expected: "apple",
		},
		{
			name:     "only unique",
			input:    "apple\napple\nbanana\n",
			opts:     Options{OnlyUnique: true},
			expected: "banana",
		},
		{
			name:     "ignore case",
			input:    "Apple\napple\nAPPLE\nBanana\n",
			opts:     Options{IgnoreCase: true},
			expected: "Apple\nBanana",
		},
		{
			name:     "skip fields and chars",
			input:    "foo bar baz\nfoo BAR baz\n",
			opts:     Options{SkipFields: 1, SkipChars: 1, IgnoreCase: true},
			expected: "foo bar baz",
		},
		{
			name:     "empty input",
			input:    "",
			opts:     Options{},
			expected: "",
		},
		{
			name:     "all unique",
			input:    "a\nb\nc\n",
			opts:     Options{},
			expected: "a\nb\nc",
		},
		{
			name:     "all duplicates",
			input:    "x\nx\nx\nx\n",
			opts:     Options{},
			expected: "x",
		},
		{
			name:     "all duplicates with count",
			input:    "x\nx\nx\nx\n",
			opts:     Options{Count: true},
			expected: "4 x",
		},
		{
			name:     "spaces only",
			input:    "   \n   \n",
			opts:     Options{},
			expected: "   ",
		},
		{
			name:     "mixed case without ignore",
			input:    "Hello\nhello\nHELLO\n",
			opts:     Options{},
			expected: "Hello\nhello\nHELLO",
		},
		{
			name:     "mixed case with ignore",
			input:    "Hello\nhello\nHELLO\n",
			opts:     Options{IgnoreCase: true},
			expected: "Hello",
		},
		{
			name:     "skip fields only",
			input:    "foo bar baz\nfoo qux baz\n",
			opts:     Options{SkipFields: 1},
			expected: "foo bar baz\nfoo qux baz",
		},
		{
			name:     "skip chars only",
			input:    "123abc\n123ABC\n",
			opts:     Options{SkipChars: 3, IgnoreCase: true},
			expected: "123abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := []string{}
			if tt.input != "" {
				lines = strings.Split(strings.TrimSuffix(tt.input, "\n"), "\n")
			}
			got := UniqueLines(lines, tt.opts)
			require.Equal(t, tt.expected, strings.Join(got, "\n"))
		})
	}
}
