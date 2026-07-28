# Test 2

Test fenced code block with `hcl`

```hcl
resource "azurerm_storage_container" "hcl" {
  bucket = "tf-test-bucket-hcl"
}
```

Test fenced code block with `tf`

```tf
resource "azurerm_storage_container" "tf" {
  bucket = "tf-test-bucket-tf"
}
```

Test block with leading whitespace

```terraform
    resource "azurerm_storage_container" "leading-space" {
  bucket = "tf-test-bucket-leading-space"
}
```

Test block with leading whitespace and line

```terraform
    
    resource "azurerm_storage_container" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line"
}
```

Test block with capital letters in resource name

```terraform
resource "azurerm_storage_container" "UpperCase" {
  bucket = "tf-test-bucket-with-uppercase"
}
```
