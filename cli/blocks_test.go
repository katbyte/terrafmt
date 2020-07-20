package cli

import (
	"strings"
	"testing"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
)

func TestCmdBlocks(t *testing.T) {
	testcases := []struct {
		name       string
		sourcefile string
		resultfile string
	}{
		{
			name:       "Go no change",
			sourcefile: "testdata/no_diffs.go",
			resultfile: "testdata/no_diffs_blocks.go.txt",
		},
		{
			name:       "Go formatting",
			sourcefile: "testdata/has_diffs.go",
			resultfile: "testdata/has_diffs_blocks.go.txt",
		},
		{
			name:       "Go fmt verbs",
			sourcefile: "testdata/fmt_compat.go",
			resultfile: "testdata/fmt_compat_blocks.go.txt",
		},
		{
			name:       "Markdown no change",
			sourcefile: "testdata/no_diffs.md",
			resultfile: "testdata/no_diffs_blocks.md.txt",
		},
		{
			name:       "Markdown formatting",
			sourcefile: "testdata/has_diffs.md",
			resultfile: "testdata/has_diffs_blocks.md.txt",
		},
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	for _, testcase := range testcases {
		data, err := afero.ReadFile(fs, testcase.resultfile)
		if err != nil {
			t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
		}
		expected := c.String(string(data))

		var outB strings.Builder
		var errB strings.Builder
		common.Log = common.CreateLogger(&errB)
		err = findBlocksInFile(fs, testcase.sourcefile, nil, &outB, &errB)
		actualStdOut := outB.String()
		actualStdErr := errB.String()

		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.sourcefile, err)
			continue
		}

		if actualStdOut != expected {
			t.Errorf("Case %q: Output does not match expected:\n%s", testcase.sourcefile, diff.Diff(actualStdOut, expected))
		}

		if actualStdErr != "" {
			t.Errorf("Case %q: Got error output:\n%s", testcase.sourcefile, actualStdErr)
		}
	}
}
