package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
)

func TestCmdFmt(t *testing.T) {
	testcases := []struct {
		name           string
		sourcefile     string
		resultfile     string
		expectErrMsg   bool
		fmtcompat      bool
		fixFinishLines bool
	}{
		{
			name:       "no change",
			sourcefile: "testdata/no_diffs.go",
			resultfile: "testdata/no_diffs.go", // No change expected
		},
		{
			name:       "formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_fmt.go",
		},
		{
			name:           "formatting, fix finish line",
			sourcefile:     "testdata/has_diffs.go",
			resultfile:     "testdata/has_diffs_fmt_fix_finish.go",
			fixFinishLines: true,
		},
		{
			name:         "fmt verbs",
			sourcefile:   "testdata/fmt_compat.go",
			resultfile:   "testdata/fmt_compat.go", // No change expected
			fmtcompat:    false,
			expectErrMsg: true,
		},
		{
			name:       "fmt verbs --fmtcompat",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_fmtcompat.go",
			fmtcompat:  true,
		},
	}

	for _, testcase := range testcases {
		inR, err := os.Open(testcase.sourcefile)
		if err != nil {
			t.Fatalf("Error reading test input file %q: %s", testcase.resultfile, err)
		}

		data, err := ioutil.ReadFile(testcase.resultfile)
		if err != nil {
			t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
		}
		expected := c.String(string(data))

		outB := bytes.NewBufferString("")
		errB := bytes.NewBufferString("")
		common.Log = common.CreateLogger(errB)
		_, err = formatFile("", testcase.fmtcompat, testcase.fixFinishLines, inR, outB, errB)
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
