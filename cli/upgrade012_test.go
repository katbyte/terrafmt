package cli

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
)

func TestCmdUpgrade012(t *testing.T) {
	testcases := []struct {
		name         string
		sourcefile   string
		resultfile   string
		noDiff       bool
		expectErrMsg bool
		fmtcompat    bool
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
			name:         "Go fmt verbs",
			sourcefile:   "testdata/fmt_compat.go",
			noDiff:       true,
			fmtcompat:    false,
			expectErrMsg: true,
		},
		{
			name:       "Go fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_upgrade012.go",
			fmtcompat:  true,
		},
	}

	for _, testcase := range testcases {
		inR, err := os.Open(testcase.sourcefile)
		if err != nil {
			t.Fatalf("Error reading test input file %q: %s", testcase.resultfile, err)
		}

		resultfile := testcase.resultfile
		if testcase.noDiff {
			resultfile = testcase.sourcefile
		}
		data, err := ioutil.ReadFile(resultfile)
		if err != nil {
			t.Fatalf("Error reading test result file %q: %s", resultfile, err)
		}
		expected := c.String(string(data))

		var outB strings.Builder
		var errB strings.Builder
		common.Log = common.CreateLogger(&errB)
		_, err = upgrade012File("", testcase.fmtcompat, inR, &outB, &errB)
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
