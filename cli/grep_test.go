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

var hclgrepTestcases = []struct {
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

  %s

  tags = {
    %[1]q = %[2]q
  }
}
`),
	},
}

func TestCmdGrepDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range hclgrepTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			opts, _, err := hclgrep.ParseArgs(testcase.hclgrep)
			if err != nil {
				t.Fatalf("parsing hclgrep args %v: %v", testcase.hclgrep, err)
			}
			br, err := hclgrepInFile(fs, log, testcase.sourcefile, false, opts, nil, &outB, &errB)
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
