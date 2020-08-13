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

func TestCmdUpgrade012StdinDefault(t *testing.T) {
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
			resultfile: "testdata/has_diffs_upgrade012.go", // This has stricter formatting than `fmt`
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			noDiff:     true,
			fmtcompat:  false,
			errMsg: []string{
				"block 1 @ stdin:8 failed to process with: cmd.Run() failed in terraform init with exit status 1:",
				"block 3 @ stdin:26 failed to process with: cmd.Run() failed in terraform init with exit status 1:",
			},
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_upgrade012.go",
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
			resultfile: "testdata/has_diffs_upgrade012.md", // This has stricter formatting than `fmt`
		},
	}

	t.Parallel()

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			inR, err := fs.Open(testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error opening test input file %q: %s", testcase.sourcefile, err)
			}

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
			_, err = upgrade012File(fs, log, "", testcase.fmtcompat, false, inR, &outB, &errB)
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

func TestCmdUpgrade012StdinVerbose(t *testing.T) {
	testcases := []struct {
		name              string
		sourcefile        string
		noDiff            bool
		lineCount         int
		updatedBlockCount int
		totalBlockCount   int
		fmtcompat         bool
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
			updatedBlockCount: 2,
			totalBlockCount:   3,
			fmtcompat:         true,
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
			lineCount:         27,
			updatedBlockCount: 3,
			totalBlockCount:   4,
		},
	}

	t.Parallel()

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			inR, err := fs.Open(testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error opening test input file %q: %s", testcase.sourcefile, err)
			}

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			_, err = upgrade012File(fs, log, "", testcase.fmtcompat, true, inR, &outB, &errB)
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
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
				t.Errorf("Unexpected summary:\nexpected %s\ngot      %s", expectedSummaryLine, summaryLine)
			}
		})
	}
}

func TestCmdUpgrade012FileDefault(t *testing.T) {
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
			resultfile: "testdata/has_diffs_upgrade012.go", // This has stricter formatting than `fmt`
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			noDiff:     true,
			fmtcompat:  false,
			errMsg: []string{
				"block 1 @ testdata/fmt_compat.go:8 failed to process with: cmd.Run() failed in terraform init with exit status 1:",
				"block 3 @ testdata/fmt_compat.go:26 failed to process with: cmd.Run() failed in terraform init with exit status 1:",
			},
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_upgrade012.go",
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
			resultfile: "testdata/has_diffs_upgrade012.md", // This has stricter formatting than `fmt`
		},
	}

	t.Parallel()

	for _, testcase := range testcases {
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
			_, err = upgrade012File(fs, log, testcase.sourcefile, testcase.fmtcompat, false, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}

			if actualStdOut != "" {
				t.Errorf("Stdout should not have output, got:\n%s", actualStdOut)
			}

			data, err = afero.ReadFile(fs, testcase.sourcefile)
			if err != nil {
				t.Fatalf("Error reading results from file %q: %s", resultfile, err)
			}
			actualContent := c.String(string(data))
			if actualContent != expected {
				t.Errorf("File does not match expected:\n%s", diff.Diff(actualContent, expected))
			}

			checkExpectedErrors(t, actualStdErr, testcase.errMsg)
		})
	}
}

func TestCmdUpgrade012FileVerbose(t *testing.T) {
	testcases := []struct {
		name              string
		sourcefile        string
		noDiff            bool
		lineCount         int
		updatedBlockCount int
		totalBlockCount   int
		fmtcompat         bool
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
			updatedBlockCount: 2,
			totalBlockCount:   3,
			fmtcompat:         true,
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
			lineCount:         27,
			updatedBlockCount: 3,
			totalBlockCount:   4,
		},
	}

	t.Parallel()

	for _, testcase := range testcases {
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
			_, err := upgrade012File(fs, log, testcase.sourcefile, testcase.fmtcompat, true, nil, &outB, &errB)
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
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
				t.Errorf("Unexpected summary:\nexpected %s\ngot      %s", expectedSummaryLine, summaryLine)
			}
		})
	}
}
