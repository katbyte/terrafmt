package test1

import (
	"fmt"
)

func testReturnSprintfSimple() string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "simple" {
  name = "tf-test-container-simple"
}
`)
}

func testReturnSprintfWithParameters(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "with-parameters" {
  name = "tf-test-container-with-parameters-%d"
}
`, randInt)
}

func testReturnSprintfWithParametersAndStringAppend(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "with-parameters-and-append" {
  name = "tf-test-container-parameters-and-append-%d"
}
`, randInt) + testReturnSprintfSimple()
}
