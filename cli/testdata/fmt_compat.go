package test3

import (
	"fmt"
)

func testNoFormattingErrors(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "no-errors" {
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
resource "azurerm_storage_container" "absolutely-nothing" {
  bucket = "tf-test-bucket-absolutely-nothing"
}
`, randInt)
}

func testExtraSpace(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "extra-space" {
  bucket    = "tf-test-bucket-extra-space-%d"

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
  bucket = "tf-test-bucket-parameters-and-append-%d"
  %[1]s     = "something"
}
`, randInt)
}

func testForExpression(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_redis_cache" "for-expression" {
  replication_group_id = %[1]q

  node_groups {
    primary_availability_zone  = azurerm_subnet.test[0].availability_zone
    replica_availability_zones = [for x in range(1, %[2]d+1) : element(azurerm_subnet.test[*].availability_zone, x)]
    replica_count              = %[2]d
  }
}
`, randInt)
}

func testFormatVerbResourceName(name string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" %[1]q {
  bucket = "tf-test-bucket-with-quotedname"
}

resource "azurerm_storage_container" "%[1]s-copy" {
  bucket = "tf-test-bucket-with-name-in-quotes"
}
`, name)
}
