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
			name: "bareverb-positional",
			block: `
%[1]s
    %[7]s
	%[77]s

%[7]d
    %[7]d

%[42]t
    %[1]t

%[7]q
    %[77]q

%[7]f
    %[77]f

%[1]g
    %[2]g

resource "resource" "test" {
	kat = "byte"
} 
`,
			expected: `
%[1]s
    %[7]s
	%[77]s

%[7]d
    %[7]d

%[42]t
    %[1]t

%[7]q
    %[77]q

%[7]f
    %[77]f

%[1]g
    %[2]g

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
mega = [%d]
    byte =   [%d]
} 
`,
			expected: `
resource "resource" "test" {
  kat  = [%s]
  mega = [%d]
  byte = [%d]
}
`,
		},

		{
			name: "assigned_array-positional",
			block: `
resource  "resource"    "test" {
	kat = [%[1]s]
mega = [%[3]d]
    byte =   [%[5]d]
} 
`,
			expected: `
resource "resource" "test" {
  kat  = [%[1]s]
  mega = [%[3]d]
  byte = [%[5]d]
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
