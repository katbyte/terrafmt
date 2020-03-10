package blocks

import "testing"

func TestBlockReaderIsFinishLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "acctest without vars",
			line:     "`)\n",
			expected: true,
		},
		{
			name:     "acctest with vars",
			line:     "`,",
			expected: true,
		},
		{
			name:     "acctest without vars and whitespaces",
			line:     "  `)\n",
			expected: true,
		},
		{
			name:     "acctest with vars and whitespaces",
			line:     "  `,",
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsFinishLine(test.line)
			if result != test.expected {
				t.Errorf("Got: \n%#v\nexpected:\n%#v\n", result, test.expected)
			}
		})
	}
}
