package format

import "testing"

func TestFmtVerbBlock(t *testing.T) {
	tests := []struct {
		name     string
		block    string
		expected string
		error    bool
	}{
		{
			name: "noverbs",
			block: `
resource  "resource"    "test" {
	kat =          "byte"
} 
`,
			expected: `
resource "resource" "test" {
  kat = "byte"
}
`,
		},

		//todo nested or forloop with letters?
		{
			name: "bareverb",
			block: `
%s
    %s
	%s

%d
    %d

%t
    %t

%q
    %q

%f
    %f

%g
    %g

resource "resource" "test" {
	kat = "byte"
} 
`,
			expected: `
%s
    %s
	%s

%d
    %d

%t
    %t

%q
    %q

%f
    %f

%g
    %g

resource "resource" "test" {
  kat = "byte"
}
`,
		},

		{
			name: "assigned_array",
			block: `
resource  "resource"    "test" {
	kat = [%s]
    byte = [%d]
} 
`,
			expected: `
resource "resource" "test" {
  kat  = [%s]
  byte = [%d]
}
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := FmtVerbBlock(test.block)
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
