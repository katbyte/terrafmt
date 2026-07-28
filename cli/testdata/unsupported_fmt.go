package test6

import (
	"fmt"
)

func testUnsupportedFmtVerbs(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "multi-verb" {
  name =    "tf-test-container"

  tags = {
    %[1]q =    %[2]q
    Test  =  "${%[5]s.name}"
    Name  =       "${%s.name}"
    byte       = "${azurerm_key_vault_certificate.test.*.id[%[2]d]}"
    Data  =    "${data.%s.name}"
  }
}
`, randInt)
}
