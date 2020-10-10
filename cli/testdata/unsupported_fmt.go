package test6

import (
	"fmt"
)

func testUnsupportedFmtVerbs(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "multi-verb" {
  bucket =    "tf-test-bucket"

  tags = {
    %[1]q =    %[2]q
    Test  =  "${%[5]s.name}"
    Name  =       "${%s.name}"
    byte       = "${aws_acm_certificate.test.*.arn[%[2]d]}"
    Data  =    "${data.%s.name}"
  }
}
`, randInt)
}
