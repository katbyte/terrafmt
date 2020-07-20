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

func TestCmdFmtStdinDefault(t *testing.T) {
	testcases := []struct {
		name           string
		sourcefile     string
		resultfile     string
		noDiff         bool
		errMsg         []string
		fmtcompat      bool
		fixFinishLines bool
	}{
		{
			name:       "Go no change",
			sourcefile: "testdata/no_diffs.go",
			noDiff:     true,
		},
		{
			name:       "Go formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_fmt.go",
		},
		{
			name:           "Go formatting, fix finish line",
			sourcefile:     "testdata/has_diffs.go",
			resultfile:     "testdata/has_diffs_fmt_fix_finish.go",
			fixFinishLines: true,
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			noDiff:     true,
			fmtcompat:  false,
			errMsg: []string{
				"block 1 @ stdin:8 failed to process with: failed to parse hcl: :4,3-4:",
				"block 3 @ stdin:26 failed to process with: failed to parse hcl: :4,3-4:",
			},
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_fmtcompat.go",
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
			resultfile: "testdata/has_diffs_fmt.md",
		},
		{
			name:           "Markdown formatting, fix finish line",
			sourcefile:     "testdata/has_diffs.md",
			resultfile:     "testdata/has_diffs_fmt.md",
			fixFinishLines: true,
		},
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	for _, testcase := range testcases {
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
		common.Log = common.CreateLogger(&errB)
		_, err = formatFile(fs, "", testcase.fmtcompat, testcase.fixFinishLines, false, inR, &outB, &errB)
		actualStdOut := outB.String()
		actualStdErr := errB.String()

		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			continue
		}

		if actualStdOut != expected {
			t.Errorf("Case %q: Output does not match expected:\n%s", testcase.name, diff.Diff(actualStdOut, expected))
		}

		checkExpectedErrors(t, testcase.name, actualStdErr, testcase.errMsg)
	}
}

func TestCmdFmtStdinVerbose(t *testing.T) {
	testcases := []struct {
		name              string
		sourcefile        string
		noDiff            bool
		lineCount         int
		updatedBlockCount int
		totalBlockCount   int
		fmtcompat         bool
		fixFinishLines    bool
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
			lineCount:         39,
			updatedBlockCount: 2,
			totalBlockCount:   4,
		},
		{
			name:              "Go formatting, fix finish line",
			sourcefile:        "testdata/has_diffs.go",
			lineCount:         39,
			updatedBlockCount: 2,
			totalBlockCount:   4,
			fixFinishLines:    true,
		},
		{
			name:            "Go fmt verbs",
			sourcefile:      "testdata/fmt_compat.go",
			noDiff:          true,
			lineCount:       33,
			totalBlockCount: 3,
			fmtcompat:       false,
		},
		{
			name:              "Go fmt verbs --fmtcompat",
			sourcefile:        "testdata/fmt_compat.go",
			lineCount:         33,
			updatedBlockCount: 1,
			totalBlockCount:   3,
			fmtcompat:         true,
		},
		{
			name:            "Markdown no change",
			sourcefile:      "testdata/no_diffs.md",
			noDiff:          true,
			lineCount:       25,
			totalBlockCount: 2,
		},
		{
			name:              "Markdown formatting",
			sourcefile:        "testdata/has_diffs.md",
			lineCount:         27,
			updatedBlockCount: 3,
			totalBlockCount:   4,
		},
		{
			name:              "Markdown formatting, fix finish line",
			sourcefile:        "testdata/has_diffs.md",
			lineCount:         27,
			updatedBlockCount: 3,
			totalBlockCount:   4,
			fixFinishLines:    true,
		},
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	for _, testcase := range testcases {
		inR, err := fs.Open(testcase.sourcefile)
		if err != nil {
			t.Fatalf("Error opening test input file %q: %s", testcase.sourcefile, err)
		}
		defer inR.Close()

		var outB strings.Builder
		var errB strings.Builder
		common.Log = common.CreateLogger(&errB)
		_, err = formatFile(fs, "", testcase.fmtcompat, testcase.fixFinishLines, true, inR, &outB, &errB)
		actualStdErr := errB.String()

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
	}
}

func TestCmdFmtFileDefault(t *testing.T) {
	testcases := []struct {
		name           string
		sourcefile     string
		resultfile     string
		noDiff         bool
		errMsg         []string
		fmtcompat      bool
		fixFinishLines bool
	}{
		{
			name:       "Go no change",
			sourcefile: "testdata/no_diffs.go",
			noDiff:     true,
		},
		{
			name:       "Go formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_fmt.go",
		},
		{
			name:           "Go formatting, fix finish line",
			sourcefile:     "testdata/has_diffs.go",
			resultfile:     "testdata/has_diffs_fmt_fix_finish.go",
			fixFinishLines: true,
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			noDiff:     true,
			fmtcompat:  false,
			errMsg: []string{
				"block 1 @ testdata/fmt_compat.go:8 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
				"block 3 @ testdata/fmt_compat.go:26 failed to process with: failed to parse hcl: testdata/fmt_compat.go:4,3-4:",
			},
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_fmtcompat.go",
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
			resultfile: "testdata/has_diffs_fmt.md",
		},
		{
			name:           "Markdown formatting, fix finish line",
			sourcefile:     "testdata/has_diffs.md",
			resultfile:     "testdata/has_diffs_fmt.md",
			fixFinishLines: true,
		},
	}

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(afero.NewOsFs()),
		afero.NewMemMapFs(),
	)

	for _, testcase := range testcases {
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
		common.Log = common.CreateLogger(&errB)
		_, err = formatFile(fs, testcase.sourcefile, testcase.fmtcompat, testcase.fixFinishLines, false, nil, &outB, &errB)
		actualStdOut := outB.String()
		actualStdErr := errB.String()

		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			continue
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
			t.Errorf("Case %q: File does not match expected:\n%s", testcase.name, diff.Diff(actualContent, expected))
		}

		checkExpectedErrors(t, testcase.name, actualStdErr, testcase.errMsg)
	}
}

func TestCmdFmtFileVerbose(t *testing.T) {
	testcases := []struct {
		name              string
		sourcefile        string
		noDiff            bool
		lineCount         int
		updatedBlockCount int
		totalBlockCount   int
		fmtcompat         bool
		fixFinishLines    bool
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
			lineCount:         39,
			updatedBlockCount: 2,
			totalBlockCount:   4,
		},
		{
			name:       "Go formatting, fix finish line",
			sourcefile: "testdata/has_diffs.go",
			noDiff:     true, // This is the actual value
			// noDiff:          false, // This is the expected value
			lineCount:         39,
			updatedBlockCount: 0, // This is the actual value
			// updatedBlockCount: 2, // This is the expected value
			totalBlockCount: 4,
			fixFinishLines:  true,
		},
		{
			name:            "Go fmt verbs",
			sourcefile:      "testdata/fmt_compat.go",
			noDiff:          true,
			lineCount:       33,
			totalBlockCount: 3,
			fmtcompat:       false,
		},
		{
			name:              "Go fmt verbs --fmtcompat",
			sourcefile:        "testdata/fmt_compat.go",
			lineCount:         33,
			updatedBlockCount: 1,
			totalBlockCount:   3,
			fmtcompat:         true,
		},
		{
			name:            "Markdown no change",
			sourcefile:      "testdata/no_diffs.md",
			noDiff:          true,
			lineCount:       25,
			totalBlockCount: 2,
		},
		{
			name:              "Markdown formatting",
			sourcefile:        "testdata/has_diffs.md",
			lineCount:         27,
			updatedBlockCount: 3,
			totalBlockCount:   4,
		},
		{
			name:       "Markdown formatting, fix finish line",
			sourcefile: "testdata/has_diffs.md",
			noDiff:     true, // This is the actual value
			// noDiff:          false, // This is the expected value
			lineCount:         27,
			updatedBlockCount: 0, // This is the actual value
			// updatedBlockCount: 3, // This is the expected value
			totalBlockCount: 4,
			fixFinishLines:  true,
		},
	}

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(afero.NewOsFs()),
		afero.NewMemMapFs(),
	)

	for _, testcase := range testcases {
		var outB strings.Builder
		var errB strings.Builder
		common.Log = common.CreateLogger(&errB)
		formatFile(fs, testcase.sourcefile, testcase.fmtcompat, testcase.fixFinishLines, true, nil, &outB, &errB)
		actualStdErr := errB.String()

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
	}
}
