package format

import (
	"log"
	"regexp"
	"strings"
)

//todo list of them?
const fmtVerbCompatibilityDelimiter = "@@_@@ TFMT"

func FmtVerbBlock(b string) (string, error) {

	// handle bare %s
	// figure out why the * doesn't match both later
	b = string(regexp.MustCompile(`(?m:^%[sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))
	b = string(regexp.MustCompile(`(?m:^[\s]*%[sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))

	// handle bare %[n]s
	b = string(regexp.MustCompile(`(?m:^%\[[\d]+\][sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))
	b = string(regexp.MustCompile(`(?m:^[\s]*%\[[\d]+\][sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))

	// handle = [%s]
	b = string(regexp.MustCompile(`(?m:\[%[sdfgtq]\])`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`))

	// handle = [%[n]s]
	b = string(regexp.MustCompile(`(?m:\[%\[[\d]+\][sdfgtq]\])`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`))

	// handle = %s/%t

	// handle = %[n]s
	// handle = [%[n]s]

	log.Print(b)

	fb, err := Block(b)
	if err != nil {
		return fb, err
	}

	//undo replace
	fb = strings.ReplaceAll(fb, "#@@_@@ TFMT:", "")

	fb = strings.ReplaceAll(fb, "[\"@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"]", "")

	return fb, nil
}
