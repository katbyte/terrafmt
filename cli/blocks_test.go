package cli

import (
	"bytes"
	"io/ioutil"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
)

func TestCmdBlocks(t *testing.T) {
	testcases := []struct {
		sourcefile string
		resultfile string
	}{
		{
			sourcefile: "testdata/no_diffs.go",
			resultfile: "testdata/no_diffs_blocks.txt",
		},
		{
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_blocks.txt",
		},
		{
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_blocks.txt",
		},
	}

	for _, testcase := range testcases {
		data, err := ioutil.ReadFile(testcase.resultfile)
		if err != nil {
			t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
		}
		expected := c.String(string(data))

		outB := bytes.NewBufferString("")
		errB := bytes.NewBufferString("")
		common.Log = common.CreateLogger(errB)
		err = findBlocksInFile(testcase.sourcefile, nil, outB, errB)
		actualOut := outB.String()
		actualErr := errB.String()

		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.sourcefile, err)
			continue
		}

		if actualOut != expected {
			t.Errorf("Case %q: Output does not match expected:\n%s", testcase.sourcefile, diff.Diff(actualOut, expected))
		}

		if actualErr != "" {
			t.Errorf("Case %q: Got error output:\n%s", testcase.sourcefile, actualErr)
		}
	}
}
