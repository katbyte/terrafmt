package test4

import (
	"fmt"
)

func testInvalidBlockName(randInt int) string {
	return fmt.Sprintf(`
rrrrrresource "azurerm_storage_container" "rrrrrrr" {
  name =    "tf-test-container"
}
`, randInt)
}

func testUnclosedBlock(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "unclosed" {
  name =    "tf-test-container"
`, randInt)
}
