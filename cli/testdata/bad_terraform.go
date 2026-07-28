package test4

import (
	"fmt"
)

func testInvalidBlockName(randInt int) string {
	return fmt.Sprintf(`
rrrrrresource "azurerm_storage_container" "rrrrrrr" {
  bucket =    "tf-test-bucket"
}
`, randInt)
}

func testUnclosedBlock(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "unclosed" {
  bucket =    "tf-test-bucket"
`, randInt)
}
