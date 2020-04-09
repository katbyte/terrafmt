# Test 2

Test fenced code block with `hcl`

```hcl
resource "aws_s3_bucket" "hcl" {
  bucket = "tf-test-bucket-hcl"
}
```

Test fenced code block with `tf`

```tf
resource "aws_s3_bucket" "tf" {
  bucket = "tf-test-bucket-tf"
}
```
