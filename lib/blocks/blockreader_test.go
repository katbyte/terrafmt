package blocks

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"gopkg.in/yaml.v2"
)

func TestBlockReaderIsFinishLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "acctest without vars",
			line:     "`)\n",
			expected: true,
		},
		{
			name:     "acctest with vars",
			line:     "`,",
			expected: true,
		},
		{
			name:     "acctest without vars and whitespaces",
			line:     "  `)\n",
			expected: true,
		},
		{
			name:     "acctest with vars and whitespaces",
			line:     "  `,",
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsFinishLine(test.line)
			if result != test.expected {
				t.Errorf("Got: \n%#v\nexpected:\n%#v\n", result, test.expected)
			}
		})
	}
}

type results struct {
	ExpectedResults []string `yaml:"expected_results"`
}

func TestBlockDetection(t *testing.T) {
	testcases := []struct {
		sourcefile string
		resultfile string
	}{
		{
			sourcefile: "testdata/test1.go",
			resultfile: "testdata/test1_results.yaml",
		},
		{
			sourcefile: "testdata/test2.markdown",
			resultfile: "testdata/test2_results.yaml",
		},
	}

	errB := bytes.NewBufferString("")
	common.Log = common.CreateLogger(errB)

	for _, testcase := range testcases {
		data, err := ioutil.ReadFile(testcase.resultfile)
		if err != nil {
			t.Fatalf("Error reading test result file %q: %s", testcase.resultfile, err)
		}
		var expectedResults results
		err = yaml.Unmarshal(data, &expectedResults)
		if err != nil {
			t.Fatalf("Error parsing test result file %q: %s", testcase.resultfile, err)
		}

		var actualBlocks []string
		br := Reader{
			ReadOnly: true,
			LineRead: ReaderIgnore,
			BlockRead: func(br *Reader, i int, b string) error {
				actualBlocks = append(actualBlocks, b)
				return nil
			},
		}
		err = br.DoTheThing(testcase.sourcefile, nil, nil)
		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.sourcefile, err)
			continue
		}

		if len(expectedResults.ExpectedResults) != len(actualBlocks) {
			t.Errorf("Case %q: expected %d blocks, got %d", testcase.sourcefile, len(expectedResults.ExpectedResults), len(actualBlocks))
			continue
		}

		for i, actual := range actualBlocks {
			expected := expectedResults.ExpectedResults[i]
			if actual != expected {
				t.Errorf("Case %q, block %d:\n%s", testcase.sourcefile, i+1, diff.Diff(expected, actual))
				continue
			}
		}

		actualErr := errB.String()
		if actualErr != "" {
			t.Errorf("Case %q: Got error output:\n%s", testcase.sourcefile, actualErr)
		}
	}
}
