package cli

import (
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/magodo/hclgrep/hclgrep"
	"github.com/spf13/afero"
)

var grepTestcases = []struct {
	name       string
	hclgrep    []string
	sourcefile string
	expected   string
}{
	{
		name:       "Go no change",
		hclgrep:    []string{"-x", `resource aws_s3_bucket simple {@*_}`},
		sourcefile: "testdata/no_diffs.go",
		expected: c.String(`
<white>#######</> <cyan>B1</><darkGray> @ #12</>
resource "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`),
	},
	{
		name:       "Go fmt verbs",
		hclgrep:    []string{"-x", `resource aws_s3_bucket no-errors {@*_}`},
		sourcefile: "testdata/fmt_compat.go",
		expected: c.String(`
<white>#######</> <cyan>B1</><darkGray> @ #18</>
resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors-%d"

#@@_@@ TFMT:  %s:TMFT @@_@@#

  tags = {
    "@@_@@ TFMT:%[1]q:TFMT @@_@@" = "@@_@@ TFMT:%[2]q:TFMT @@_@@"
  }
}
`),
	},
}

func TestCmdGrepDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range grepTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			m := hclgrep.Matcher{}
			cmds, _, err := m.ParseCmds(testcase.hclgrep)
			if err != nil {
				t.Fatalf("parsing hclgrep args %v: %v", testcase.hclgrep, err)
			}
			br, err := grepInFile(fs, log, testcase.sourcefile, false, cmds, m, nil, &outB, &errB)
			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}
			if br.ErrorBlocks != 0 {
				t.Errorf("Expected no block errors, got %d", br.ErrorBlocks)
			}
			actualContent := outB.String()
			if actualContent != testcase.expected {
				t.Errorf("Case %q: File does not match expected: ('-' actual, '+' expected)\n%s", testcase.name, diff.Diff(actualContent, testcase.expected))
			}
		})
	}
}
