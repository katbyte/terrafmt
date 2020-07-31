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

func TestCmdDiffDefault(t *testing.T) {
	testcases := []struct {
		name       string
		sourcefile string
		resultfile string
		noDiff     bool
		errMsg     []string
		fmtcompat  bool
	}{
		{
			name:       "Go no change",
			sourcefile: "testdata/no_diffs.go",
			noDiff:     true,
		},
		{
			name:       "Go formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_diff.go.txt",
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_diff_nofmtcompat.go.txt",
			fmtcompat:  false,
			errMsg: []string{
				"block 1 @ testdata/fmt_compat.go:8 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
				"block 3 @ testdata/fmt_compat.go:26 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
			},
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_diff_fmtcompat.go.txt",
			fmtcompat:  true,
		},
		{
			name:       "Markdown no change",
			sourcefile: "testdata/no_diffs.md",
			noDiff:     true,
		},
		{
			name:       "Markdown formatting",
			sourcefile: "testdata/has_diffs.md",
			resultfile: "testdata/has_diffs_diff.md.txt",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
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
			common.Log = common.CreateLogger(&errB)
			_, _, err := diffFile(fs, testcase.sourcefile, testcase.fmtcompat, false, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}

			if actualStdOut != expected {
				t.Errorf("Output does not match expected:\n%s", diff.Diff(actualStdOut, expected))
			}

			checkExpectedErrors(t, actualStdErr, testcase.errMsg)
		})
	}
}

func TestCmdDiffVerbose(t *testing.T) {
	testcases := []struct {
		name                  string
		sourcefile            string
		noDiff                bool
		lineCount             int
		unformattedBlockCount int
		totalBlockCount       int
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
			lineCount:             39,
			unformattedBlockCount: 2,
			totalBlockCount:       4,
		},
		{
			name:            "Go fmt verbs",
			sourcefile:      "testdata/fmt_compat.go",
			noDiff:          true, // The only diff is in the block with the parsing error
			lineCount:       33,
			totalBlockCount: 3,
			fmtcompat:       false,
		},
		{
			name:                  "Go fmt verbs --fmtcompat",
			sourcefile:            "testdata/fmt_compat.go",
			lineCount:             33,
			unformattedBlockCount: 1,
			totalBlockCount:       3,
			fmtcompat:             true,
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
			lineCount:             27,
			unformattedBlockCount: 3,
			totalBlockCount:       4,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			var outB strings.Builder
			var errB strings.Builder
			common.Log = common.CreateLogger(&errB)
			_, _, err := diffFile(fs, testcase.sourcefile, testcase.fmtcompat, true, nil, &outB, &errB)
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
