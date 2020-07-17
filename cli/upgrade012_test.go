package cli

import (
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
)

func TestCmdUpgrade012(t *testing.T) {
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

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	for _, testcase := range testcases {
		inR, err := fs.Open(testcase.sourcefile)
		if err != nil {
			t.Fatalf("Error reading test input file %q: %s", testcase.resultfile, err)
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
		common.Log = common.CreateLogger(&errB)
		_, err = upgrade012File(fs, "", testcase.fmtcompat, inR, &outB, &errB)
		actualOut := outB.String()
		actualErr := errB.String()

		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			continue
		}

		if actualOut != expected {
			t.Errorf("Case %q: Output does not match expected:\n%s", testcase.name, diff.Diff(actualOut, expected))
		}

		checkExpectedErrors(t, testcase.name, actualErr, testcase.errMsg)
	}
}
