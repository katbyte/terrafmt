package test3

import (
	"fmt"
)

func testNoFormattingErrors(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "no-errors" {
  name = "tf-test-container-no-errors-%d"

  %s

  tags = {
    %[1]q = %[2]q
  }
}
`, randInt)
}

func testNoErrorsOrFmtVerbs(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "absolutely-nothing" {
  name = "tf-test-container-absolutely-nothing"
}
`, randInt)
}

func testExtraSpace(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "extra-space" {
  name    = "tf-test-container-extra-space-%d"

  %s

  tags = {
    %[1]q    = %[2]q
  }
}
`, randInt) + testReturnSprintfSimple()
}

func testFormatVerbParameter(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "with-parameters-and-append" {
  name = "tf-test-container-parameters-and-append-%d"
  %[1]s     = "something"
}
`, randInt)
}

func testForExpression(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_redis_cache" "for-expression" {
  replication_group_id = %[1]q

  node_groups {
    primary_address_prefixes  = azurerm_subnet.test[0].address_prefixes
    replica_address_prefixess = [for x in range(1, %[2]d+1) : element(azurerm_subnet.test[*].address_prefixes, x)]
    replica_count              = %[2]d
  }
}
`, randInt)
}

func testFormatVerbResourceName(name string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" %[1]q {
  name = "tf-test-container-with-quotedname"
}

resource "azurerm_storage_container" "%[1]s-copy" {
  name = "tf-test-container-with-name-in-quotes"
}
`, name)
}
