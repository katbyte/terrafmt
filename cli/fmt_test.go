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

var fmtTestcases = []struct {
	name              string
	sourcefile        string
	resultfile        string
	noDiff            bool
	errMsg            []string
	fmtcompat         bool
	fixFinishLines    bool
	lineCount         int
	updatedBlockCount int
	totalBlockCount   int
}{
	{
		name:            "Go no change",
		sourcefile:      "testdata/no_diffs.go",
		noDiff:          true,
		lineCount:       29,
		totalBlockCount: 3,
	},
	{
		name:              "Go formatting",
		sourcefile:        "testdata/has_diffs.go",
		resultfile:        "testdata/has_diffs_fmt.go",
		lineCount:         86,
		updatedBlockCount: 4,
		totalBlockCount:   6,
	},
	{
		name:              "Go formatting, fix finish line",
		sourcefile:        "testdata/has_diffs.go",
		resultfile:        "testdata/has_diffs_fmt_fix_finish.go",
		fixFinishLines:    true,
		lineCount:         86,
		updatedBlockCount: 5,
		totalBlockCount:   6,
	},
	{
		name:       "Go fmt verbs",
		sourcefile: "testdata/fmt_compat.go",
		noDiff:     true,
		fmtcompat:  false,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: %s:4,3-4:",
			"block 3 @ %s:30 failed to process with: failed to parse hcl: %s:4,3-4:",
			"block 4 @ %s:44 failed to process with: failed to parse hcl: %s:3,3-4:",
			"block 5 @ %s:53 failed to process with: failed to parse hcl: %s:2,26-27:",
		},
		lineCount:       64,
		totalBlockCount: 5,
	},
	{
		name:              "Go fmt verbs --fmtcompat",
		sourcefile:        "testdata/fmt_compat.go",
		resultfile:        "testdata/fmt_compat_fmtcompat.go",
		fmtcompat:         true,
		lineCount:         64,
		updatedBlockCount: 3,
		totalBlockCount:   5,
	},
	{
		name:       "Go bad terraform",
		sourcefile: "testdata/bad_terraform.go",
		resultfile: "testdata/bad_terraform_fmt.go",
		errMsg: []string{
			"block 2 @ %s:16 failed to process with: failed to parse hcl: %s:1,37-38: Unclosed configuration block; There is no closing brace for this block before the end of the file. This may be caused by incorrect brace nesting elsewhere in this file.\\nresource \\",
		},
		lineCount:         20,
		updatedBlockCount: 1,
		totalBlockCount:   2,
	},
	{
		name:       "Go unsupported format verbs",
		sourcefile: "testdata/unsupported_fmt.go",
		noDiff:     true,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: %s:5,5-6:",
		},
		lineCount:       21,
		totalBlockCount: 1,
	},
	{
		name:       "Go unsupported format verbs --fmtcompat",
		sourcefile: "testdata/unsupported_fmt.go",
		noDiff:     true,
		fmtcompat:  true,
		errMsg: []string{
			"block 1 @ %s:8 failed to process with: failed to parse hcl: %s:6,17-18:",
		},
		lineCount:       21,
		totalBlockCount: 1,
	},
	{
		name:            "Markdown no change",
		sourcefile:      "testdata/no_diffs.md",
		noDiff:          true,
		lineCount:       25,
		totalBlockCount: 3,
	},
	{
		name:              "Markdown formatting",
		sourcefile:        "testdata/has_diffs.md",
		resultfile:        "testdata/has_diffs_fmt.md",
		lineCount:         33,
		updatedBlockCount: 4,
		totalBlockCount:   5,
	},
	{
		name:              "Markdown formatting, fix finish line",
		sourcefile:        "testdata/has_diffs.md",
		resultfile:        "testdata/has_diffs_fmt.md",
		fixFinishLines:    true,
		lineCount:         33,
		updatedBlockCount: 4,
		totalBlockCount:   5,
	},
}

func TestCmdFmtStdinDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range fmtTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			inR, err := fs.Open(testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error opening test input file %q: %s", testcase.sourcefile, err)
			}
			defer inR.Close()

			resultfile := testcase.resultfile
			if testcase.noDiff {
				resultfile = testcase.sourcefile
			}
			data, err := afero.ReadFile(fs, resultfile)
			if err != nil {
				t.Fatalf("Error reading test result file %q: %s", resultfile, err)
			}
			expected := c.String(string(data))

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			br, err := formatFile(fs, log, "", testcase.fmtcompat, testcase.fixFinishLines, false, inR, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			if len(testcase.errMsg) != br.ErrorBlocks {
				t.Errorf("Expected %d block errors, got %d", len(testcase.errMsg), br.ErrorBlocks)
			}

			if actualStdOut != expected {
				t.Errorf("Case %q: Output does not match expected: ('-' actual, '+' expected)\n%s", testcase.name, diff.Diff(actualStdOut, expected))
			}

			errMsg := []string{}
			for _, msg := range testcase.errMsg {
				errMsg = append(errMsg, fmt.Sprintf(msg, "stdin", ""))
			}
			checkExpectedErrors(t, actualStdErr, errMsg)
		})
	}
}

func TestCmdFmtStdinVerbose(t *testing.T) {
	t.Parallel()

	for _, testcase := range fmtTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			inR, err := fs.Open(testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error opening test input file %q: %s", testcase.sourcefile, err)
			}
			defer inR.Close()

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			_, err = formatFile(fs, log, "", testcase.fmtcompat, testcase.fixFinishLines, true, inR, &outB, &errB)
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			filenameColor := "lightMagenta"
			if testcase.noDiff {
				filenameColor = "magenta"
			}
			expectedSummaryLine := c.String(fmt.Sprintf(
				"<%s>%s</>: <cyan>%d</> lines & formatted <yellow>%d</>/<yellow>%d</> blocks!",
				filenameColor,
				"stdin",
				testcase.lineCount,
				testcase.updatedBlockCount,
				testcase.totalBlockCount,
			))

			trimmedStdErr := strings.TrimSpace(actualStdErr)
			lines := strings.Split(trimmedStdErr, "\n")
			summaryLine := lines[len(lines)-1]
			if summaryLine != expectedSummaryLine {
				t.Errorf("Case %q: Unexpected summary:\nexpected %s\ngot      %s", testcase.name, expectedSummaryLine, summaryLine)
			}
		})
	}
}

func TestCmdFmtFileDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range fmtTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewCopyOnWriteFs(
				afero.NewReadOnlyFs(afero.NewOsFs()),
				afero.NewMemMapFs(),
			)

			resultfile := testcase.resultfile
			if testcase.noDiff {
				resultfile = testcase.sourcefile
			}
			data, err := afero.ReadFile(fs, resultfile)
			if err != nil {
				t.Fatalf("Error reading test result file %q: %s", resultfile, err)
			}
			expected := c.String(string(data))

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			br, err := formatFile(fs, log, testcase.sourcefile, testcase.fmtcompat, testcase.fixFinishLines, false, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			if actualStdOut != "" {
				t.Errorf("Case %q: Stdout should not have output, got:\n%s", testcase.name, actualStdOut)
			}

			data, err = afero.ReadFile(fs, testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error reading results from file %q: %s", resultfile, err)
			}
			actualContent := c.String(string(data))
			if actualContent != expected {
				t.Errorf("Case %q: File does not match expected: ('-' actual, '+' expected)\n%s", testcase.name, diff.Diff(actualContent, expected))
			}

			if len(testcase.errMsg) != br.ErrorBlocks {
				t.Errorf("Expected %d block errors, got %d", len(testcase.errMsg), br.ErrorBlocks)
			}

			errMsg := []string{}
			for _, msg := range testcase.errMsg {
				errMsg = append(errMsg, fmt.Sprintf(msg, testcase.sourcefile, testcase.sourcefile))
			}
			checkExpectedErrors(t, actualStdErr, errMsg)
		})
	}
}

func TestCmdFmtFileVerbose(t *testing.T) {
	t.Parallel()

	for _, testcase := range fmtTestcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewCopyOnWriteFs(
				afero.NewReadOnlyFs(afero.NewOsFs()),
				afero.NewMemMapFs(),
			)

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			_, err := formatFile(fs, log, testcase.sourcefile, testcase.fmtcompat, testcase.fixFinishLines, true, nil, &outB, &errB)
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			filenameColor := "lightMagenta"
			if testcase.noDiff {
				filenameColor = "magenta"
			}
			expectedSummaryLine := c.String(fmt.Sprintf(
				"<%s>%s</>: <cyan>%d</> lines & formatted <yellow>%d</>/<yellow>%d</> blocks!",
				filenameColor,
				testcase.sourcefile,
				testcase.lineCount,
				testcase.updatedBlockCount,
				testcase.totalBlockCount,
			))

			trimmedStdErr := strings.TrimSpace(actualStdErr)
			lines := strings.Split(trimmedStdErr, "\n")
			summaryLine := lines[len(lines)-1]
			if summaryLine != expectedSummaryLine {
				t.Errorf("Case %q: Unexpected summary:\nexpected %s\ngot      %s", testcase.name, expectedSummaryLine, summaryLine)
			}
		})
	}
}
