package test4

import (
	"fmt"
)

func testInvalidBlockName(randInt int) string {
	return fmt.Sprintf(`
rrrrrresource "aws_s3_bucket" "rrrrrrr" {
  bucket =    "tf-test-bucket"
}
`, randInt)
}

func testUnclosedBlock(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "unclosed" {
  bucket =    "tf-test-bucket"
`, randInt)
}
