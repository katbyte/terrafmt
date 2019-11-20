package upgrade012

import (
	"regexp"
	"strings"
)

//todo list of them?
const fmtVerbCompatibilityDelimiter = "@@_@@ TFMT"

func FmtVerbBlock(b string) (string, error) {

	// NOTE: the order of these replacements matter

	// handle bare %s
	// figure out why the * doesn't match both later
	b = string(regexp.MustCompile(`(?m:^%[sdfgtqv]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))
	b = string(regexp.MustCompile(`(?m:^[\s]*%[sdfgtqv]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))

	// handle bare %[n]s
	b = string(regexp.MustCompile(`(?m:^%\[[\d]+\][sdfgtqv]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))
	b = string(regexp.MustCompile(`(?m:^[\s]*%\[[\d]+\][sdfgtqv]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))

	// handle = [%s]
	b = string(regexp.MustCompile(`(?m:\[%[sdfgtqv]\]$)`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`))

	// handle = [%[n]s]
	b = string(regexp.MustCompile(`(?m:\[%\[[\d]+\][sdfgtqv]\]$)`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`))

	// handle = %s/%t
	b = string(regexp.MustCompile(`(?m:%[sdfgtqv]$)`).ReplaceAllString(b, `"@@_@@ TFMT:$0:TFMT @@_@@"`))

	// handle = %[n]s
	b = string(regexp.MustCompile(`(?m:%\[[\d]+\][sdfgtqv]$)`).ReplaceAllString(b, `"@@_@@ TFMT:$0:TFMT @@_@@"`))

	// handle = "%s"
	b = string(regexp.MustCompile(`(?m:\"%([sdfgtqv])\")`).ReplaceAllString(b, `"TFMT__${1}__TFMT"`))

	// handle = .%s.
	b = string(regexp.MustCompile(`(?m:\.%([sdfgtqv])\.)`).ReplaceAllString(b, `.TFMT__${1}__TFMT.`))

	fb, err := Block(b)
	if err != nil {
		return fb, err
	}

	//undo replace
	fb = strings.ReplaceAll(fb, "#@@_@@ TFMT:", "")

	fb = strings.ReplaceAll(fb, "[\"@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"]", "")

	// order here matters, replace the ones with [], then do the ones without
	fb = strings.ReplaceAll(fb, "\"@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"", "")

	fb = strings.ReplaceAll(fb, "TFMT__", "%")
	fb = strings.ReplaceAll(fb, "__TFMT", "")

	return fb, nil
}
