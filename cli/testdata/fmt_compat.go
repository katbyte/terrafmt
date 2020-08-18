package test3

import (
	"fmt"
)

func testNoFormattingErrors(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors-%d"

  %s

  tags = {
    %[1]q = %[2]q
  }
}
`, randInt)
}

func testNoErrorsOrFmtVerbs(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "absolutely-nothing" {
  bucket = "tf-test-bucket-absolutely-nothing"
}
`, randInt)
}

func testExtraSpace(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "extra-space" {
  bucket    = "tf-test-bucket-extra-space-%d"

  %s

  tags = {
    %[1]q    = %[2]q
  }
}
`, randInt) + testReturnSprintfSimple()
}
