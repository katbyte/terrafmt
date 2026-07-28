# Test 2

Test fenced code block with `hcl`

```hcl
resource "azurerm_storage_container" "hcl" {
  name = "tf-test-container-hcl"
}
```

Test fenced code block with `tf`

```tf
resource "azurerm_storage_container" "tf" {
  name = "tf-test-container-tf"
}
```

Test block with leading whitespace

```terraform
    resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space"
}
```

Test block with leading whitespace and line

```terraform
    
    resource "azurerm_storage_container" "leading-space-and-line" {
  name = "tf-test-container-leading-space-and-line"
}
```

Test block with capital letters in resource name

```terraform
resource "azurerm_storage_container" "UpperCase" {
  name = "tf-test-container-with-uppercase"
}
```
