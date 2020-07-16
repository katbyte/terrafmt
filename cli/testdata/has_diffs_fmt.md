# Has Diffs

```hcl
resource "aws_s3_bucket" "extra-lines" {

  bucket = "tf-test-bucket-extra-lines"
}
```

```hcl
resource "aws_s3_bucket" "no-errors" {
  bucket = "tf-test-bucket-no-errors-%d"
}
```

```hcl
resource "aws_s3_bucket" "extra-space" {
  bucket = "tf-test-bucket-extra-space-%d"
}
```

```hcl
resource "aws_s3_bucket" "end-line" {
  bucket = "tf-test-bucket-end-line-%d"
}

```
