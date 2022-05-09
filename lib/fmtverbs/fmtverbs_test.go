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
  kat  = base64encode(%[1]s)
  byte = md5(data.source.%[1]s.id)
}
`,
			expected: `
resource "resource" "test" {
  kat  = base64encode(TFMTFNPARAM_s)
  byte = md5(data.source.TFMTKTKTTFMTs.id)
  kat  = base64encode(TFMTFNPARAM_1s)
  byte = md5(data.source.TFMTKTKTTFMT_1s.id)
}
`,
		},
		{
			name: "012 syntax properties",
			block: `
resource "resource" "test" {
  s = data.source.%s.id
  d = data.source.%d.id
  t = data.source.%t.id
}
`,
			expected: `
resource "resource" "test" {
  s = data.source.TFMTKTKTTFMTs.id
  d = data.source.TFMTKTKTTFMTd.id
  t = data.source.TFMTKTKTTFMTt.id
}
`,
		},
		{
			name: "012 syntax properties positional",
			block: `
resource "resource" "test" {
  s = data.source.%[1]s.id
  d = data.source.%[2]d.id
  t = data.source.%[5]t.id
}
`,
			expected: `
resource "resource" "test" {
  s = data.source.TFMTKTKTTFMT_1s.id
  d = data.source.TFMTKTKTTFMT_2d.id
  t = data.source.TFMTKTKTTFMT_5t.id
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
    Ωq = "@@_@@ TFMT:%q:TFMT @@_@@"
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
    Ω_1_q = "@@_@@ TFMT:%[2]q:TFMT @@_@@"
  }
}
`,
		},
		{
			// No change expected
			name: "looks like 012",
			block: `
resource "resource" "test" {
  kat  = "%s.example.com"
  byte = "%[1]s.example.com"
}
`,
			expected: `
resource "resource" "test" {
  kat  = "%s.example.com"
  byte = "%[1]s.example.com"
}
`,
		},
		{
			name: "old-style",
			block: `
resource "resource" "test" {
  kat  = "${%s.name}"
  byte = "${%[5]s.name}"
}

resource "resource" "test" {
  kat  = "${aws_acm_certificate.test.*.arn[%d]}"
  byte = "${aws_acm_certificate.test.*.arn[%[2]d]}"
}
`,
			expected: `
resource "resource" "test" {
  kat  = "${%s.name}"
  byte = "${%[5]s.name}"
}

resource "resource" "test" {
  kat  = "${aws_acm_certificate.test.*.arn[0/*@@_@@ TFMT:%d:TFMT @@_@@*/]}"
  byte = "${aws_acm_certificate.test.*.arn[0/*@@_@@ TFMT:%[2]d:TFMT @@_@@*/]}"
}
`,
		},
		{
			name: "conditional expression",
			block: `
resource "aws_dynamodb_table" "test" {
  name = %[1]q

  ttl {
    attribute_name = %[2]t ? "TestTTL" : ""
    enabled        = %[2]t
  }
}
`,
			expected: `
resource "aws_dynamodb_table" "test" {
  name = "@@_@@ TFMT:%[1]q:TFMT @@_@@"

  ttl {
    attribute_name = true/*@@_@@ TFMT:%[2]t:TFMT @@_@@*/ ? "TestTTL" : ""
    enabled        = "@@_@@ TFMT:%[2]t:TFMT @@_@@"
  }
}
`,
		},
		{
			name: "verb in index",
			block: `
resource "resource" "test" {
  attr = aws_acm_certificate.test[%[2]d].arn
  attr = "${aws_acm_certificate.test.*.arn[%[2]d]}"
  attr = aws_acm_certificate.test[%d].arn
  attr = "${aws_acm_certificate.test.*.arn[%d]}"
}
`,
			expected: `
resource "resource" "test" {
  attr = aws_acm_certificate.test["@@_@@ TFMT:[%[2]d]:TFMT @@_@@"].arn
  attr = "${aws_acm_certificate.test.*.arn[0/*@@_@@ TFMT:%[2]d:TFMT @@_@@*/]}"
  attr = aws_acm_certificate.test["@@_@@ TFMT:[%d]:TFMT @@_@@"].arn
  attr = "${aws_acm_certificate.test.*.arn[0/*@@_@@ TFMT:%d:TFMT @@_@@*/]}"
}
`,
		},
		{
			name: "verb as parameter name",
			block: `
resource "resource" "test1" {
  %s = %q
}

resource "resource" "test2" {
  %[1]s = %q
}

resource "resource" "test3" {
  %s = %[3]q
}

resource "resource" "test4" {
  %[6]s = %[2]q
}

resource "resource" "test5" {
  %s = {
    %[3]q
  }
}

resource "resource" "test6" {
  %[4]s = {
    %[2]q
  }
}
`,
			expected: `
resource "resource" "test1" {
  Ωs = "@@_@@ TFMT:%q:TFMT @@_@@"
}

resource "resource" "test2" {
  Ω_1_s = "@@_@@ TFMT:%q:TFMT @@_@@"
}

resource "resource" "test3" {
  Ωs = "@@_@@ TFMT:%[3]q:TFMT @@_@@"
}

resource "resource" "test4" {
  Ω_6_s = "@@_@@ TFMT:%[2]q:TFMT @@_@@"
}

resource "resource" "test5" {
  Ωs = {
#@@_@@ TFMT:    %[3]q:TMFT @@_@@#
  }
}

resource "resource" "test6" {
  Ω_4_s = {
#@@_@@ TFMT:    %[2]q:TMFT @@_@@#
  }
}
`,
		},
		{
			name: "verb in for expression",
			block: `
resource "resource" "test" {
  attr = [for x in range(1, %d+1) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(1, %[2]d+1) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(%d, 3) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(%[1]d, 3) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(%d, %d) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(%[1]d, %[2]d) : element(aws_subnet.test[*].availability_zone, x)]
}
`,
			expected: `
resource "resource" "test" {
  attr = [for x in range(1, TFMTFNPARAM_d+1) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(1, TFMTFNPARAM_2d+1) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(TFMTFNPARAM_d, 3) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(TFMTFNPARAM_1d, 3) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(TFMTFNPARAM_d, TFMTFNPARAM_d) : element(aws_subnet.test[*].availability_zone, x)]
  attr = [for x in range(TFMTFNPARAM_1d, TFMTFNPARAM_2d) : element(aws_subnet.test[*].availability_zone, x)]
}
`,
		},
		{
			name: "verb in resource name",
			block: `
resource "resource" "%s" {
  kat = resource.test-%s.byte
}

resource "resource" "test-%[1]s" {
  kat = resource.%s-test.byte
}

resource "resource" "%s-test" {
  kat = resource.%s.byte
}

data "data_source" %[1]q {
}

resource "resource" %q {
  kat = resource.%[1]s.byte
}
`,
			expected: `
resource "resource" "TFMTRESNAME_s" {
  kat = resource.test-TFMTKTKTTFMTs.byte
}

resource "resource" "test-TFMTRESNAME_1s" {
  kat = resource.TFMTKTKTTFMTs-test.byte
}

resource "resource" "TFMTRESNAME_s-test" {
  kat = resource.TFMTKTKTTFMTs.byte
}

data "data_source" "TFMTRESNAME_1q" {
}

resource "resource" "TFMTRESNAME_q" {
  kat = resource.TFMTKTKTTFMT_1s.byte
}
`,
		},
		{
			name: "provider meta-argument",
			block: `
resource "resource" "test" {
  provider = %s
}

resource "resource" "test2" {
  provider = %[1]s
}
`,
			expected: `
resource "resource" "test" {
  provider = tfmtprovider.PROVIDER
}

resource "resource" "test2" {
  provider = tfmtprovider.PROVIDER_1
}
`,
		},
		{
			name: "count meta-argument",
			block: `
resource "resource" "test" {
  count = %d
}

resource "resource" "test2" {
  count = %[2]d
}
`,
			expected: `
resource "resource" "test" {
  count = var.tfmtcount
}

resource "resource" "test2" {
  count = var.tfmtcount_2
}
`,
		},
	}

	t.Parallel()

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

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
