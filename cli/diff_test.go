package cli

import (
	"fmt"
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
)

var diffTestcases = []struct {
	name                  string
	sourcefile            string
	resultfile            string
	noDiff                bool
	lineCount             int
	unformattedBlockCount int
	totalBlockCount       int
	errMsg                []string
	fmtcompat             bool
}{
	{
		name:            "Go no change",
		sourcefile:      "testdata/no_diffs.go",
		noDiff:          true,
		lineCount:       29,
		totalBlockCount: 3,
	},
	{
		name:                  "Go formatting",
		sourcefile:            "testdata/has_diffs.go",
		resultfile:            "testdata/has_diffs_diff.go.txt",
		lineCount:             86,
		unformattedBlockCount: 4,
		totalBlockCount:       6,
	},
	{
		name:       "Go fmt verbs",
		sourcefile: "testdata/fmt_compat.go",
		resultfile: "testdata/fmt_compat_diff_nofmtcompat.go.txt",
		fmtcompat:  false,
		noDiff:     true,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
			"block 3 @ %s:30 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
			"block 4 @ %s:44 failed to process with: failed to parse hcl: testdata/fmt_compat.go:3,3-4:",
			"block 5 @ %s:53 failed to process with: failed to parse hcl: testdata/fmt_compat.go:2,26-27:",
			"block 6 @ %s:67 failed to process with: failed to parse hcl: testdata/fmt_compat.go:1,26-27:",
		},
		lineCount:       76,
		totalBlockCount: 6,
	},
	{
		name:                  "Go fmt verbs --fmtcompat",
		sourcefile:            "testdata/fmt_compat.go",
		resultfile:            "testdata/fmt_compat_diff_fmtcompat.go.txt",
		fmtcompat:             true,
		lineCount:             76,
		unformattedBlockCount: 3,
		totalBlockCount:       6,
	},
	{
		name:       "Go bad terraform",
		sourcefile: "testdata/bad_terraform.go",
		resultfile: "testdata/bad_terraform_diff.go.txt",
		errMsg: []string{
			"block 2 @ %s:16 failed to process with: failed to parse hcl: testdata/bad_terraform.go:1,37-38: Unclosed configuration block; There is no closing brace for this block before the end of the file. This may be caused by incorrect brace nesting elsewhere in this file.\\nresource \\",
		},
		lineCount:             20,
		unformattedBlockCount: 1,
		totalBlockCount:       2,
	},
	{
		name:       "Go unsupported format verbs",
		sourcefile: "testdata/unsupported_fmt.go",
		noDiff:     true,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: testdata/unsupported_fmt.go:5,5-6:",
		},
		lineCount:             21,
		unformattedBlockCount: 0,
		totalBlockCount:       1,
	},
	{
		name:       "Go unsupported format verbs --fmtcompat",
		sourcefile: "testdata/unsupported_fmt.go",
		noDiff:     true,
		fmtcompat:  true,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: testdata/unsupported_fmt.go:6,17-18:",
		},
		lineCount:             21,
		unformattedBlockCount: 0,
		totalBlockCount:       1,
	},
	{
		name:            "Markdown no change",
		sourcefile:      "testdata/no_diffs.md",
		noDiff:          true,
		lineCount:       25,
		totalBlockCount: 3,
	},
	{
		name:                  "Markdown formatting",
		sourcefile:            "testdata/has_diffs.md",
		resultfile:            "testdata/has_diffs_diff.md.txt",
		lineCount:             33,
		unformattedBlockCount: 4,
		totalBlockCount:       5,
	},
}

func TestCmdDiffDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range diffTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			expected := ""
			if !testcase.noDiff {
				data, err := afero.ReadFile(fs, testcase.resultfile)
				if err != nil {
					t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
				}
				expected = c.String(string(data))
			}

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			br, hasDiff, err := diffFile(fs, log, testcase.sourcefile, testcase.fmtcompat, false, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}

			actualNoDiff := !hasDiff
			if testcase.noDiff && !actualNoDiff {
				t.Errorf("Expected no diff, but got one")
			} else if !testcase.noDiff && actualNoDiff {
				t.Errorf(("Expected diff, but did not get one"))
			}

			if len(testcase.errMsg) != br.ErrorBlocks {
				t.Errorf("Expected %d block errors, got %d", len(testcase.errMsg), br.ErrorBlocks)
			}

			if actualStdOut != expected {
				t.Errorf("Output does not match expected: ('-' actual, '+' expected)\n%s", diff.Diff(actualStdOut, expected))
			}

			errMsg := []string{}
			for _, msg := range testcase.errMsg {
				errMsg = append(errMsg, fmt.Sprintf(msg, testcase.sourcefile))
			}
			checkExpectedErrors(t, actualStdErr, errMsg)
		})
	}
}

func TestCmdDiffVerbose(t *testing.T) {
	t.Parallel()

	for _, testcase := range diffTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			_, _, err := diffFile(fs, log, testcase.sourcefile, testcase.fmtcompat, true, nil, &outB, &errB)
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}

			filenameColor := "lightMagenta"
			if testcase.noDiff {
				filenameColor = "magenta"
			}
			expectedSummaryLine := c.String(fmt.Sprintf(
				"<%s>%s</>: <cyan>%d</> lines & <yellow>%d</>/<yellow>%d</> blocks need formatting.",
				filenameColor,
				testcase.sourcefile,
				testcase.lineCount,
				testcase.unformattedBlockCount,
				testcase.totalBlockCount,
			))

			trimmedStdErr := strings.TrimSpace(actualStdErr)
			lines := strings.Split(trimmedStdErr, "\n")
			summaryLine := lines[len(lines)-1]
			if summaryLine != expectedSummaryLine {
				t.Errorf("Unexpected summary:\nexpected %s\ngot      %s", expectedSummaryLine, summaryLine)
			}
		})
	}
}
