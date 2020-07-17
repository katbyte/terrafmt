package cli

import (
	"io/ioutil"
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
)

func TestCmdDiff(t *testing.T) {
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
		expected := ""
		if !testcase.noDiff {
			data, err := ioutil.ReadFile(testcase.resultfile)
			if err != nil {
				t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
			}
			expected = c.String(string(data))
		}

		var outB strings.Builder
		var errB strings.Builder
		common.Log = common.CreateLogger(&errB)
		_, _, err := diffFile(testcase.sourcefile, testcase.fmtcompat, nil, &outB, &errB)
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
