package test1

import (
	"fmt"
)

func testReturnSprintfSimple() string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`)
}

func testReturnStringSimple() string {
	return `
resource "aws_s3_bucket" "simple2" {
  bucket = "tf-test-bucket-simple2"
}
`
}

func testReturnSprintfWithParameters(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "with-parameters" {
  bucket = "tf-test-bucket-with-parameters-%d"
}
`, randInt)
}

func testReturnSprintfWithParametersAndStringAppend(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "with-parameters-and-append" {
  bucket = "tf-test-bucket-parameters-and-append-%d"
}
`, randInt) + testReturnSprintfSimple()
}
