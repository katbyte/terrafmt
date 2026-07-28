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

func testReturnStringSimple() string {
	return `
resource "azurerm_storage_container" "simple2" {
  name = "tf-test-container-simple2"
}
`
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

const testConst = `
resource "azurerm_storage_container" "const" {
  name = "tf-test-container-const"
}
`

func testComposed(randInt int) string {
	return testReturnSprintfWithParameters(randInt) + fmt.Sprintf(`
resource "azurerm_storage_container" "composed" {
  name = "tf-test-container-composed-%d"
}
`, randInt)
}

func testDataSource() string {
	return fmt.Sprintf(`
data "azurerm_storage_container" "simple" {
  name = "tf-test-container-simple"
}
`)
}

func testLeadingWhiteSpace(randInt int) string {
	return fmt.Sprintf(`
    resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space-%d"
}
`, randInt)
}

func testLeadingWhiteSpaceAndLine(randInt int) string {
	return fmt.Sprintf(`
    
    resource "azurerm_storage_container" "leading-space-and-line" {
  name = "tf-test-container-leading-space-and-line-%d"
}
`, randInt)
}

func testFormatVerbResourceName(name string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "%s" {
  name = "tf-test-container-with-quotedname"
}
`, name)
}

func testFormatUpperCase(name string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_container" "UpperCase" {
  name = "tf-test-container-with-uppercase"
}
`)
}

func notTerraformSimpleString() string {
	fmt.Sprintf("%d: bad create: \n%#v\n%#v", i, cm, tc.Create)
}

func notTerraformXML() string {
	return `<DescribeAccountAttributesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
	<requestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</requestId>
	<accountAttributeSet>
	  <item>
		<attributeName>supported-platforms</attributeName>
		<attributeValueSet>
		  <item>
			<attributeValue>VPC</attributeValue>
		  </item>
		  <item>
			<attributeValue>EC2</attributeValue>
		  </item>
		</attributeValueSet>
	  </item>
	</accountAttributeSet>
  </DescribeAccountAttributesResponse>`
}
