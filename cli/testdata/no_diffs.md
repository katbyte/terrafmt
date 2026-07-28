# No Diffs

```terraform
resource "azurerm_storage_container" "one" {
  bucket = "tf-test-bucket-one"
}
```

```hcl
resource "azurerm_storage_container" "two" {
  bucket = "tf-test-bucket-two"
}
```

```tf
resource "azurerm_storage_container" "three" {
  bucket = "tf-test-bucket-three"
}
```

```
resource "azurerm_storage_container" "four" {
  bucket = "tf-test-bucket-four"
}
```
