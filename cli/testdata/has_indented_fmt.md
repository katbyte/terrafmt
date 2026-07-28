# indented fences

1. A block that needs formatting:

    ```hcl
    resource "aws_s3_bucket" "one" {
      bucket        = "tf-test-bucket-one"
      force_destroy = true
    }
    ```

2. A block that is already formatted:

    ```tf
    resource "aws_s3_bucket" "two" {
      bucket = "tf-test-bucket-two"
    }
    ```
