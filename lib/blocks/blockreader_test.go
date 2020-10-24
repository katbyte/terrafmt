package blocks

import (
	"bytes"
	"testing"

	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
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

func TestBlockDetection(t *testing.T) {
	testcases := []struct {
		sourcefile     string
		expectedBlocks []string
	}{
		{
			sourcefile: "testdata/test1.go",
			expectedBlocks: []string{
				`resource "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`,
				`resource "aws_s3_bucket" "with-parameters" {
  bucket = "tf-test-bucket-with-parameters-%d"
}
`,
				`resource "aws_s3_bucket" "with-parameters-and-append" {
  bucket = "tf-test-bucket-parameters-and-append-%d"
}
`,
				`resource "aws_s3_bucket" "const" {
  bucket = "tf-test-bucket-const"
}
`,
				`resource "aws_s3_bucket" "composed" {
  bucket = "tf-test-bucket-composed-%d"
}
`,
				`data "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`,
				`    resource "aws_s3_bucket" "leading-space" {
  bucket = "tf-test-bucket-leading-space-%d"
}
`,
				`    
    resource "aws_s3_bucket" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line-%d"
}
`,
			},
		},
		{
			sourcefile: "testdata/test2.markdown",
			expectedBlocks: []string{
				`resource "aws_s3_bucket" "hcl" {
  bucket = "tf-test-bucket-hcl"
}
`,
				`resource "aws_s3_bucket" "tf" {
  bucket = "tf-test-bucket-tf"
}
`,
				`    resource "aws_s3_bucket" "leading-space" {
  bucket = "tf-test-bucket-leading-space"
}
`,
				`    
    resource "aws_s3_bucket" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line"
}
`,
			},
		},
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	errB := bytes.NewBufferString("")
	log := common.CreateLogger(errB)

	for _, testcase := range testcases {
		var actualBlocks []string
		br := Reader{
			Log:      log,
			ReadOnly: true,
			LineRead: ReaderIgnore,
			BlockRead: func(br *Reader, i int, b string) error {
				actualBlocks = append(actualBlocks, b)
				return nil
			},
		}
		err := br.DoTheThingNew(fs, testcase.sourcefile, nil, nil)
		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.sourcefile, err)
			continue
		}

		if len(testcase.expectedBlocks) != len(actualBlocks) {
			t.Errorf("Case %q: expected %d blocks, got %d", testcase.sourcefile, len(testcase.expectedBlocks), len(actualBlocks))
			continue
		}

		for i, actual := range actualBlocks {
			expected := testcase.expectedBlocks[i]
			if actual != expected {
				t.Errorf("Case %q, block %d: ('-' actual, '+' expected)\n%s", testcase.sourcefile, i+1, diff.Diff(actual, expected))
				continue
			}
		}

		actualErr := errB.String()
		if actualErr != "" {
			t.Errorf("Case %q: Got error output:\n%s", testcase.sourcefile, actualErr)
		}
	}
}

func TestLooksLikeTerraform(t *testing.T) {
	testcases := []struct {
		text     string
		expected bool
	}{
		{
			text: `
resource "aws_s3_bucket" "simple-resource" {
  bucket = "tf-test-bucket-simple"
}`,
			expected: true,
		},
		{
			text: `
data "aws_s3_bucket" "simple-data" {
  bucket = "tf-test-bucket-simple"
}`,
			expected: true,
		},
		{
			text: `
variable "name" {
  type = string
}`,
			expected: true,
		},
		{
			text: `
output "arn" {
  value = aws_s3_bucket.simple-resource.arn
}`,
			expected: true,
		},
		{
			text:     "%d: bad create: \n%#v\n%#v",
			expected: false,
		},
		{
			text: `<DescribeAccountAttributesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
  <requestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</requestId>
  <accountAttributeSet>
	  <item>
	    <attributeName>supported-platforms</attributeName>
	    <attributeValueSet>
	      <item>
	        <attributeValue>VPC</attributeValue>
	      </item>
	      <item>
	        <attributeValue>EC2</attributeValue>
	      </item>
	    </attributeValueSet>
	  </item>
  </accountAttributeSet>
</DescribeAccountAttributesResponse>`,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := looksLikeTerraform(testcase.text)

		if testcase.expected && !actual {
			t.Errorf("Expected match, but not identified as Terraform:\n%s", testcase.text)
		} else if !testcase.expected && actual {
			t.Errorf("Expected no match, but was identified as Terraform:\n%s", testcase.text)
		}
	}
}
