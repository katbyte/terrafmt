# Has Diffs

```hcl
resource "azurerm_storage_container" "extra-lines" {

  name = "tf-test-container-extra-lines"
}
```

```hcl
resource "azurerm_storage_container" "no-errors" {
  name = "tf-test-container-no-errors"
}
```

```hcl
resource "azurerm_storage_container" "extra-space" {
  name = "tf-test-container-extra-space"
}
```

```hcl
resource "azurerm_storage_container" "end-line" {
  name = "tf-test-container-end-line"
}

```

```hcl
resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space"
}
```
