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
		name         string
		sourcefile   string
		resultfile   string
		noDiff       bool
		expectErrMsg bool
		fmtcompat    bool
	}{
		{
			name:       "no change",
			sourcefile: "testdata/no_diffs.go",
			noDiff:     true,
		},
		{
			name:       "formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_diff.txt",
		},
		{
			name:         "fmt verbs",
			sourcefile:   "testdata/fmt_compat.go",
			resultfile:   "testdata/fmt_compat_diff_nofmtcompat.txt",
			fmtcompat:    false,
			expectErrMsg: true,
		},
		{
			name:       "fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_diff_fmtcompat.txt",
			fmtcompat:  true,
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

		if testcase.expectErrMsg {
			if strings.TrimSpace(actualErr) == "" {
				t.Errorf("Case %q: Expected error output but got none", testcase.name)
			}
		} else {
			if actualErr != "" {
				t.Errorf("Case %q: Got error output:\n%s", testcase.name, actualErr)
			}
		}
	}
}
