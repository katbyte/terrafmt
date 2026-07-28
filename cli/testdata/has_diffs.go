package test2

import (
	"fmt"
)

func testExtraLines() string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "extra-lines" {
  
  name = "tf-test-container-extra-lines"
}
`)
}

// This is included to verify blocks with diffs and no diffs in the same file
func testNoFormattingErrors(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "no-errors" {
  name = "tf-test-container-no-errors-%d"
}
`, randInt)
}

func testExtraSpace(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "extra-space" {
  name    = "tf-test-container-extra-space-%d"
}
`, randInt) + testReturnSprintfSimple()
}

func testFinishLineWhiteSpace(randInt int) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "end-line" {
  name = "tf-test-container-end-line-%d"
}
  `, randInt)
}

func testNoPadding(randInt int) string {
	return fmt.Sprintf(`resource "azurerm_lb_backend_address_pool" "test" {
  name = "%s"
  port = 443
  protocol = "HTTPS"
  vpc_id = "${azurerm_virtual_network.test.id}"

  deregistration_delay = 200

  stickiness {
    type = "lb_cookie"
    cookie_duration = 10000
  }

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }

  tags = {
    TestName = "TestAccAWSALBTargetGroup_basic"
  }
}

resource "azurerm_virtual_network" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-alb-target-group-basic"
  }
}`, targetGroupName)
}

func testLeadingWhiteSpace(randInt int) string {
	return fmt.Sprintf(`
    resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space-%d"
}
`, randInt)
}
