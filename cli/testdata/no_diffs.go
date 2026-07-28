package test1

import (
	"fmt"
)

func testReturnSprintfSimple() string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "simple" {
  bucket = "tf-test-bucket-simple"
}
`)
}

func testReturnSprintfWithParameters(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "with-parameters" {
  bucket = "tf-test-bucket-with-parameters-%d"
}
`, randInt)
}

func testReturnSprintfWithParametersAndStringAppend(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "with-parameters-and-append" {
  bucket = "tf-test-bucket-parameters-and-append-%d"
}
`, randInt) + testReturnSprintfSimple()
}
