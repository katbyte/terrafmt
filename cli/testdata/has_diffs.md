# Has Diffs

```hcl
resource "azurerm_storage_container" "extra-lines" {
  
  bucket = "tf-test-bucket-extra-lines"
}
```

```hcl
resource "azurerm_storage_container" "no-errors" {
  bucket = "tf-test-bucket-no-errors"
}
```

```hcl
resource "azurerm_storage_container" "extra-space" {
  bucket    = "tf-test-bucket-extra-space"
}
```

```hcl
resource "azurerm_storage_container" "end-line" {
  bucket = "tf-test-bucket-end-line"
}
  
```

```hcl
     resource "azurerm_storage_container" "leading-space" {
  bucket = "tf-test-bucket-leading-space"
}
```
