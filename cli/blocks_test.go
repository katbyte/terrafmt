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

type block struct {
	endLine int
	text    string
}

var testcases = []struct {
	name           string
	sourcefile     string
	lineCount      int
	expectedBlocks []block
}{
	{
		name:       "Go no change",
		sourcefile: "testdata/no_diffs.go",
		lineCount:  29,
		expectedBlocks: []block{
			{
				endLine: 12,
				text: `resource "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}`,
			},
			{
				endLine: 20,
				text: `resource "aws_s3_bucket" "with-parameters" {
  bucket = "tf-test-bucket-with-parameters-%d"
}`,
			},
			{
				endLine: 28,
				text: `resource "aws_s3_bucket" "with-parameters-and-append" {
  bucket = "tf-test-bucket-parameters-and-append-%d"
}`,
			},
		},
	},
	{
		name:       "Go formatting",
		sourcefile: "testdata/has_diffs.go",
		lineCount:  39,
		expectedBlocks: []block{
			{
				endLine: 13,
				text: `resource "aws_s3_bucket" "extra-lines" {
  
  bucket = "tf-test-bucket-extra-lines"
}`,
			},
			{
				endLine: 22,
				text: `resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors-%d"
}`,
			},
			{
				endLine: 30,
				text: `resource "aws_s3_bucket" "extra-space" {
  bucket    = "tf-test-bucket-extra-space-%d"
}`,
			},
			{
				endLine: 38,
				text: `resource "aws_s3_bucket" "end-line" {
  bucket = "tf-test-bucket-end-line-%d"
}`,
			},
		},
	},
	{
		name:       "Go fmt verbs",
		sourcefile: "testdata/fmt_compat.go",
		lineCount:  41,
		expectedBlocks: []block{
			{
				endLine: 18,
				text: `resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors-%d"

  %s

  tags = {
    %[1]q = %[2]q
  }
}`,
			},
			{
				endLine: 26,
				text: `resource "aws_s3_bucket" "absolutely-nothing" {
  bucket = "tf-test-bucket-absolutely-nothing"
}`,
			},
			{
				endLine: 40,
				text: `resource "aws_s3_bucket" "extra-space" {
  bucket    = "tf-test-bucket-extra-space-%d"

  %s

  tags = {
    %[1]q    = %[2]q
  }
}`,
			},
		},
	},
	{
		name:       "Markdown no change",
		sourcefile: "testdata/no_diffs.md",
		lineCount:  25,
		expectedBlocks: []block{
			{
				endLine: 7,
				text: `resource "aws_s3_bucket" "one" {
  bucket = "tf-test-bucket-one"
}`,
			},
			{
				endLine: 13,
				text: `resource "aws_s3_bucket" "two" {
  bucket = "tf-test-bucket-two"
}`,
			},
			{
				endLine: 19,
				text: `resource "aws_s3_bucket" "three" {
  bucket = "tf-test-bucket-three"
}`,
			},
		},
	},
	{
		name:       "Markdown formatting",
		sourcefile: "testdata/has_diffs.md",
		lineCount:  27,
		expectedBlocks: []block{
			{
				endLine: 8,
				text: `resource "aws_s3_bucket" "extra-lines" {
  
  bucket = "tf-test-bucket-extra-lines"
}`,
			},
			{
				endLine: 14,
				text: `resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors"
}`,
			},
			{
				endLine: 20,
				text: `resource "aws_s3_bucket" "extra-space" {
  bucket    = "tf-test-bucket-extra-space"
}`,
			},
			{
				endLine: 27,
				text: `resource "aws_s3_bucket" "end-line" {
  bucket = "tf-test-bucket-end-line"
}
  `,
			},
		},
	},
}

func TestCmdBlocksDefault(t *testing.T) {
	t.Parallel()

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			expectedBuilder := strings.Builder{}
			for i, block := range testcase.expectedBlocks {
				fmt.Fprint(&expectedBuilder, c.Sprintf("\n<white>#######</> <cyan>B%d</><darkGray> @ #%d</>\n", i+1, block.endLine))
				fmt.Fprint(&expectedBuilder, block.text, "\n")
			}
			expected := expectedBuilder.String()

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			err := findBlocksInFile(fs, log, testcase.sourcefile, false, false, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Got an error when none was expected: %v", err)
			}

			if actualStdOut != expected {
				t.Errorf("Output does not match expected: ('-' actual, '+' expected)\n%s", diff.Diff(actualStdOut, expected))
			}

			if actualStdErr != "" {
				t.Errorf("Got error output:\n%s", actualStdErr)
			}
		})
	}
}

func TestCmdBlocksVerbose(t *testing.T) {
	t.Parallel()

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			err := findBlocksInFile(fs, log, testcase.sourcefile, true, false, nil, &outB, &errB)
			actualStdErr := errB.String()
			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			expectedSummaryLine := c.String(fmt.Sprintf("Finished processing <cyan>%d</> lines <yellow>%d</> blocks!", testcase.lineCount, len(testcase.expectedBlocks)))

			summaryLine := strings.TrimSpace(actualStdErr)
			if summaryLine != expectedSummaryLine {
				t.Errorf("Case %q: Unexpected summary:\nexpected %s\ngot      %s", testcase.name, expectedSummaryLine, summaryLine)
			}
		})
	}
}

func TestCmdBlocksZeroTerminated(t *testing.T) {
	t.Parallel()

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewReadOnlyFs(afero.NewOsFs())

			expectedBuilder := strings.Builder{}
			for _, block := range testcase.expectedBlocks {
				fmt.Fprint(&expectedBuilder, block.text, "\n\x00")
			}
			expected := expectedBuilder.String()

			var outB strings.Builder
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			err := findBlocksInFile(fs, log, testcase.sourcefile, false, true, nil, &outB, &errB)
			actualStdOut := outB.String()
			actualStdErr := errB.String()

			if err != nil {
				t.Fatalf("Case %q: Got an error when none was expected: %v", testcase.name, err)
			}

			if actualStdOut != expected {
				t.Errorf("Case %q: Output does not match expected: ('-' actual, '+' expected)\n%s", testcase.name, diff.Diff(actualStdOut, expected))
			}

			if actualStdErr != "" {
				t.Errorf("Case %q: Got error output:\n%s", testcase.name, actualStdErr)
			}
		})
	}
}
