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

Test block with leading whitespace

```terraform
    resource "aws_s3_bucket" "leading-space" {
  bucket = "tf-test-bucket-leading-space"
}
```

Test block with leading whitespace and line

```terraform
    
    resource "aws_s3_bucket" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line"
}
```

Test block with capital letters in resource name

```terraform
resource "aws_s3_bucket" "UpperCase" {
  bucket = "tf-test-bucket-with-uppercase"
}
```
