package format

import "testing"

func TestBlock(t *testing.T) {
	tests := []struct {
		name     string
		block    string
		expected string
		error    bool
	}{
		{
			name:     "empty",
			block:    "",
			expected: "\n",
		},
		{
			name: "oneline",
			block: `
resource   "resource"    "test" {} 
`,
			expected: `
resource "resource" "test" {}
`,
		},
		{
			name: "basic",
			block: `
resource   "resource"    "test" {
	kay = "tee"
} 
`,
			expected: `
resource "resource" "test" {
  kay = "tee"
}
`,
		},
		{
			name: "whitespace",
			block: `
resource "resource" "test" {
noindent = "test"
	tabbed = "test"
         spaces {
              k = "t"
              kkkk = "t"

              kt = "kt"
     }
}  
`,
			expected: `
resource "resource" "test" {
  noindent = "test"
  tabbed   = "test"
  spaces {
    k    = "t"
    kkkk = "t"

    kt = "kt"
  }
}
`,
		},
		{
			name: "invalid",
			block: `
Hi there i am going to fail... =C
`,
			expected: ``,
			error:    true,
		},
		{
			name: "fmtVerbs",
			block: `
resource "resource" "test" {
%s
}  
`,
			expected: ``,
			error:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Block(test.block)
			if err != nil && !test.error {
				t.Fatalf("Got an error when none was expected: %v", err)
			}
			if err == nil && test.error {
				t.Errorf("Expected an error and none was generated")
			}
			if result != test.expected {
				t.Errorf("Got: \n%#v\nexpected:\n%#v\n", result, test.expected)
			}
		})
	}
}
