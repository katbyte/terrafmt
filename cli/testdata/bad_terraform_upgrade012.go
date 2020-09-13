package test4

import (
	"fmt"
)

func testInvalidBlockName(randInt int) string {
	return fmt.Sprintf(`
# TF-UPGRADE-TODO: Block type was not recognized, so this block and its contents were not automatically upgraded.
rrrrrresource "aws_s3_bucket" "rrrrrrr" {
  bucket = "tf-test-bucket"
}
`, randInt)
}

func testUnclosedBlock(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "unclosed" {
  bucket =    "tf-test-bucket"
`, randInt)
}
