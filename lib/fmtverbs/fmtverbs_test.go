package fmtverbs

import (
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestFmtVerbBlock(t *testing.T) {
	tests := []struct {
		name     string
		block    string
		expected string
	}{
		{
			name: "noverbs",
			block: `
resource  "resource"   "test" {
  kat =          "byte"
} 
`,
			expected: `
resource  "resource"   "test" {
  kat =          "byte"
} 
`,
		},
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
`,
			expected: `
#@@_@@ TFMT:%s:TMFT @@_@@#
#@@_@@ TFMT:    %s:TMFT @@_@@#
#@@_@@ TFMT:	%s:TMFT @@_@@#

#@@_@@ TFMT:%d:TMFT @@_@@#
#@@_@@ TFMT:    %d:TMFT @@_@@#

#@@_@@ TFMT:%t:TMFT @@_@@#
#@@_@@ TFMT:    %t:TMFT @@_@@#

#@@_@@ TFMT:%q:TMFT @@_@@#
#@@_@@ TFMT:    %q:TMFT @@_@@#

#@@_@@ TFMT:%f:TMFT @@_@@#
#@@_@@ TFMT:    %f:TMFT @@_@@#

#@@_@@ TFMT:%g:TMFT @@_@@#
#@@_@@ TFMT:    %g:TMFT @@_@@#
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
`,
			expected: `
#@@_@@ TFMT:%[1]s:TMFT @@_@@#
#@@_@@ TFMT:    %[7]s:TMFT @@_@@#
#@@_@@ TFMT:	%[77]s:TMFT @@_@@#

#@@_@@ TFMT:%[7]d:TMFT @@_@@#
#@@_@@ TFMT:    %[7]d:TMFT @@_@@#

#@@_@@ TFMT:%[42]t:TMFT @@_@@#
#@@_@@ TFMT:    %[1]t:TMFT @@_@@#

#@@_@@ TFMT:%[7]q:TMFT @@_@@#
#@@_@@ TFMT:    %[77]q:TMFT @@_@@#

#@@_@@ TFMT:%[7]f:TMFT @@_@@#
#@@_@@ TFMT:    %[77]f:TMFT @@_@@#

#@@_@@ TFMT:%[1]g:TMFT @@_@@#
#@@_@@ TFMT:    %[2]g:TMFT @@_@@#
`,
		},

		{
			name: "assigned-array",
			block: `
resource "resource" "test" {
  kat  = [%s]
  mega = [%d]
  byte = [%d]
  size = ["%s"]
}
`,
			expected: `
resource "resource" "test" {
  kat  = ["@@_@@ TFMT:[%s]:TFMT @@_@@"]
  mega = ["@@_@@ TFMT:[%d]:TFMT @@_@@"]
  byte = ["@@_@@ TFMT:[%d]:TFMT @@_@@"]
  size = ["%s"]
}
`,
		},
		{
			name: "assigned-array-positional",
			block: `
resource "resource" "test" {
  kat  = [%[1]s]
  mega = [%[1]d]
  byte = [%[1]d]
  size = ["%[1]s"]
}
`,
			expected: `
resource "resource" "test" {
  kat  = ["@@_@@ TFMT:[%[1]s]:TFMT @@_@@"]
  mega = ["@@_@@ TFMT:[%[1]d]:TFMT @@_@@"]
  byte = ["@@_@@ TFMT:[%[1]d]:TFMT @@_@@"]
  size = ["%[1]s"]
}
`,
		},

		{
			name: "assigned-verb",
			block: `
resource "resource" "test" {
  kat  = %s
  mega = %d
  byte = %d
  size = "%s"
}
`,
			expected: `
resource "resource" "test" {
  kat  = "@@_@@ TFMT:%s:TFMT @@_@@"
  mega = "@@_@@ TFMT:%d:TFMT @@_@@"
  byte = "@@_@@ TFMT:%d:TFMT @@_@@"
  size = "%s"
}
`,
		},
		{
			name: "assigned-verb-not-standalone",
			block: `
resource "resource" "test" {
  kat  = %s.id
  byte = %d.id
}
`,
			expected: `
resource "resource" "test" {
  kat  = "@@_@@ TFMT:%s.id:TFMT @@_@@"
  byte = "@@_@@ TFMT:%d.id:TFMT @@_@@"
}
`,
		},
		{
			name: "assigned-positional",
			block: `
resource "resource" "test" {
  kat  = %[1]s
  mega = %[22]d
  byte = %[333]d
  size = "%[42]s"
}
`,
			expected: `
resource "resource" "test" {
  kat  = "@@_@@ TFMT:%[1]s:TFMT @@_@@"
  mega = "@@_@@ TFMT:%[22]d:TFMT @@_@@"
  byte = "@@_@@ TFMT:%[333]d:TFMT @@_@@"
  size = "%[42]s"
}
`,
		},
		{
			name: "assigned-positional-not-standalone",
			block: `
resource "resource" "test" {
  kat  = %[1]s.id
  byte = %[1]d.id
}
`,
			expected: `
resource "resource" "test" {
  kat  = "@@_@@ TFMT:%[1]s.id:TFMT @@_@@"
  byte = "@@_@@ TFMT:%[1]d.id:TFMT @@_@@"
}
`,
		},
		{
			name: "function call",
			block: `
resource "resource" "test" {
  kat  = base64encode(%s)
  byte = md5(data.source.%s.id)
}
`,
			expected: `
resource "resource" "test" {
  kat  = base64encode(TFFMTKTBRACKETPERCENTs)
  byte = md5(data.source.TFMTKTKTTFMTs.id)
}
`,
		},
		{
			name: "012 syntax properties",
			block: `
resource "resource" "test" {
  s = data.source.%s.id
  d = data.source.%d.id
  f = data.source.%f.id
  g = data.source.%g.id
  t = data.source.%t.id
  q = data.source.%q.id
}
`,
			expected: `
resource "resource" "test" {
  s = data.source.TFMTKTKTTFMTs.id
  d = data.source.TFMTKTKTTFMTd.id
  f = data.source.TFMTKTKTTFMTf.id
  g = data.source.TFMTKTKTTFMTg.id
  t = data.source.TFMTKTKTTFMTt.id
  q = data.source.TFMTKTKTTFMTq.id
}
`,
		},
		{
			name: "012 syntax properties positional",
			block: `
resource "resource" "test" {
  s = data.source.%[1]s.id
  d = data.source.%[2]d.id
  f = data.source.%[3]f.id
  g = data.source.%[4]g.id
  t = data.source.%[5]t.id
  q = data.source.%[6]q.id
}
`,
			expected: `
resource "resource" "test" {
  s = data.source.TFMTKTKTTFMT_1s.id
  d = data.source.TFMTKTKTTFMT_2d.id
  f = data.source.TFMTKTKTTFMT_3f.id
  g = data.source.TFMTKTKTTFMT_4g.id
  t = data.source.TFMTKTKTTFMT_5t.id
  q = data.source.TFMTKTKTTFMT_6q.id
}
`,
		},
		{
			name: "multiple verbs",
			block: `
resource "resource" "test" {
  kat  = [%s, %s]
  mega = [%d, %d]
  byte = [%d, %d]
  size = ["%s", "%s"]

  tags = {
    %q = %q
  }
}

resource "resource" "test" {
  kat  = [%s, %s,%s]
  mega = [%d, %d,%d]
  byte = [%d, %d,%d]
  size = ["%s", "%s","%s"]
}
`,
			expected: `
resource "resource" "test" {
  kat  = ["@@_@@ TFMT:[%s, %s]:TFMT @@_@@"]
  mega = ["@@_@@ TFMT:[%d, %d]:TFMT @@_@@"]
  byte = ["@@_@@ TFMT:[%d, %d]:TFMT @@_@@"]
  size = ["%s", "%s"]

  tags = {
    "@@_@@ TFMT:%q:TFMT @@_@@" = "@@_@@ TFMT:%q:TFMT @@_@@"
  }
}

resource "resource" "test" {
  kat  = ["@@_@@ TFMT:[%s, %s,%s]:TFMT @@_@@"]
  mega = ["@@_@@ TFMT:[%d, %d,%d]:TFMT @@_@@"]
  byte = ["@@_@@ TFMT:[%d, %d,%d]:TFMT @@_@@"]
  size = ["%s", "%s","%s"]
}
`,
		},
		{
			name: "multiple verbs positional",
			block: `
resource "resource" "test" {
  kat  = [%[1]s, %[2]s,%[3]s]
  mega = [%[1]d, %[3]d,%[2]d]
  byte = [%[1]d, %[3]d,%[2]d]
  size = ["%[1]s", "%[2]s","%[3]s"]

  tags = {
    %[1]q = %[2]q
  }
}
`,
			expected: `
resource "resource" "test" {
  kat  = ["@@_@@ TFMT:[%[1]s, %[2]s,%[3]s]:TFMT @@_@@"]
  mega = ["@@_@@ TFMT:[%[1]d, %[3]d,%[2]d]:TFMT @@_@@"]
  byte = ["@@_@@ TFMT:[%[1]d, %[3]d,%[2]d]:TFMT @@_@@"]
  size = ["%[1]s", "%[2]s","%[3]s"]

  tags = {
    "@@_@@ TFMT:%[1]q:TFMT @@_@@" = "@@_@@ TFMT:%[2]q:TFMT @@_@@"
  }
}
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Escape(test.block)
			if result != test.expected {
				t.Fatalf("Unexpected escaped result: ('-' actual, '+' expected)\n%s\n", diff.Diff(result, test.expected))
			}

			roundtrip := Unscape(result)
			if roundtrip != test.block {
				t.Fatalf("Did not roundtrip: ('-' actual, '+' expected)\n%s\n", diff.Diff(roundtrip, test.block))
			}
		})
	}
}
