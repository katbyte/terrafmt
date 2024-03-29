package fmtverbs

import (
	"regexp"
	"strings"
)

func Escape(b string) string {
	// NOTE: the order of these replacements matter

	// conditional expression: = %t ? ...
	b = regexp.MustCompile(`(=\s*)(%\[[\d+]\]t)(\s\?)`).ReplaceAllString(b, `${1}true/*@@_@@ TFMT:${2}:TFMT @@_@@*/${3}`)

	// resource name contains %[sdtfg]
	b = regexp.MustCompile(`([resource|data]\s+"[-a-zA-Z0-9_]+"\s+"[-a-zA-Z0-9_]*)%(?:\[([\d]+)\])?([sdtfg])`).ReplaceAllString(b, `${1}TFMTRESNAME_${2}${3}`)

	// resource name %q
	b = regexp.MustCompile(`([resource|data]\s+"[-a-zA-Z0-9_]+"\s+)%(?:\[([\d]+)\])?q`).ReplaceAllString(b, `${1}"TFMTRESNAME_${2}q"`)

	// %s - whole line
	// figure out why the * doesn't match both later
	b = regexp.MustCompile(`(?m:^%(\.[0-9])?[sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0:TMFT @@_@@#`)
	b = regexp.MustCompile(`(?m:^[ \t]*%(\.[0-9])?[sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0:TMFT @@_@@#`)

	// provider meta-argument
	// The provider name must be in lowercase
	b = regexp.MustCompile(`(\bprovider\s+=\s+)%s`).ReplaceAllString(b, `${1}tfmtprovider.PROVIDER`)
	b = regexp.MustCompile(`(\bprovider\s+=\s+)+%\[(\d+)\]s`).ReplaceAllString(b, `${1}tfmtprovider.PROVIDER_${2}`)

	// count meta-argument
	b = regexp.MustCompile(`(\bcount\s+=\s+)%([ds])`).ReplaceAllString(b, `${1}1 # tfmtcount_${2}`)
	b = regexp.MustCompile(`(\bcount\s+=\s+)+%(\[(\d+)\][ds])`).ReplaceAllString(b, `${1}1 # tfmtcount_${2}`)

	// %[n]s
	b = regexp.MustCompile(`(?m:^%(\.[0-9])?\[[\d]+\][sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0:TMFT @@_@@#`)
	b = regexp.MustCompile(`(?m:^[ \t]*%(\.[0-9])?\[[\d]+\][sdfgtq]$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0:TMFT @@_@@#`)

	// %s =
	b = regexp.MustCompile(`(?m)^([ \t]*)%((\.[0-9])?[sdfgtq])`).ReplaceAllString(b, `${1}Ω${2}`)

	// %[n]s =
	b = regexp.MustCompile(`(?m)^([ \t]*)%(\.[0-9])?\[([\d]+)\]([sdfgtq])`).ReplaceAllString(b, `${1}Ω_${3}_${4}`)

	// = "${...[%([n])d]}"
	b = regexp.MustCompile(`(?m)("\${.*\[)(%(?:\.[0-9])?(?:\[[\d]+\])?d)(\]}")$`).ReplaceAllString(b, `${1}0/*@@_@@ TFMT:$2:TFMT @@_@@*/$3`)

	// = [%s(, %s)]
	b = regexp.MustCompile(`(?m:\[(%(\.[0-9])?[sdfgtq](,\s*)?)+\])`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`)

	// = [%[n]s(, %[n]s)]
	b = regexp.MustCompile(`(?m:\[(%(\.[0-9])?\[[\d]+\][sdfgtq](,\s*)?)+\])`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`)

	//  .12 - something.%s.prop
	b = regexp.MustCompile(`\.%([sdt])`).ReplaceAllString(b, `.TFMTKTKTTFMT${1}`)

	//  .12 - something.%[n]s.prop
	b = regexp.MustCompile(`\.%\[(\d+)\]([sdt])`).ReplaceAllString(b, `.TFMTKTKTTFMT_${1}${2}`)

	//  .12 - something.text%s.prop
	b = regexp.MustCompile(`\.([-a-zA-Z0-9_]+)%([sdt])`).ReplaceAllString(b, `.${1}TFMTKTKTTFMT${2}`)

	//  .12 - something.text%[n]s.prop
	b = regexp.MustCompile(`\.([-a-zA-Z0-9_]+)%\[(\d+)\]([sdt])`).ReplaceAllString(b, `.${1}TFMTKTKTTFMT_${2}${3}`)

	// = %s
	b = regexp.MustCompile(`(?m:%(\.[0-9])?[sdfgtq](\.[a-z_]+)*$)`).ReplaceAllString(b, `"@@_@@ TFMT:$0:TFMT @@_@@"`)

	// = %[n]s
	b = regexp.MustCompile(`(?m:%(\.[0-9])?\[[\d]+\][sdfgtq](\.[a-z_]+)*$)`).ReplaceAllString(b, `"@@_@@ TFMT:$0:TFMT @@_@@"`)

	// function(..., %s, ...)
	// function(..., %[n]s, ...)
	b = regexp.MustCompile(`\(.*%(\[(\d+)\])?[sdfgtq].*?\)`).ReplaceAllStringFunc(b, func(v string) string {
		result := regexp.MustCompile(`%([sdfgtq])`).ReplaceAllString(v, `TFMTFNPARAM_$1`)
		result = regexp.MustCompile(`%\[(\d+)\]([sdfgtq])`).ReplaceAllString(result, `TFMTFNPARAM_$1$2`)
		return result
	})

	return b
}

func Unscape(fb string) string {
	// NOTE: the order of these replacements matter

	// undo replace
	fb = regexp.MustCompile(`[ ]*#@@_@@ TFMT:`).ReplaceAllString(fb, ``)
	fb = strings.ReplaceAll(fb, "#@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TMFT @@_@@#", "")

	fb = strings.ReplaceAll(fb, "[\"@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"]", "")

	fb = strings.ReplaceAll(fb, "0/*@@_@@ TFMT:", "")
	fb = regexp.MustCompile(`true\s*/\*@@_@@ TFMT:`).ReplaceAllString(fb, ``)
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@*/", "")

	// order here matters, replace the ones with [], then do the ones without
	fb = strings.ReplaceAll(fb, "\"@@_@@ TFMT:", "")
	fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"", "")

	// .12 - something.%[n]s.prop
	fb = regexp.MustCompile(`TFMTKTKTTFMT_(\d+)`).ReplaceAllString(fb, `%[$1]`)

	// .12 - something.%s.prop
	fb = strings.ReplaceAll(fb, "TFMTKTKTTFMT", "%")

	// function(..., %[n]s, ...)
	fb = regexp.MustCompile(`TFMTFNPARAM_(\d+)([sdfgtq])`).ReplaceAllString(fb, `%[$1]$2`)

	// function(..., %s, ...)
	fb = strings.ReplaceAll(fb, "TFMTFNPARAM_", "%")

	// %[n]s =
	fb = regexp.MustCompile(`Ω_(\d+)_([sdfgtq])`).ReplaceAllString(fb, `%[${1}]${2}`)

	// %s =
	fb = regexp.MustCompile(`Ω([sdfgtq])`).ReplaceAllString(fb, `%${1}`)

	// resource name %[n]s
	fb = regexp.MustCompile(`TFMTRESNAME_(\d+)([sdtfg])`).ReplaceAllString(fb, `%[${1}]${2}`)

	// resource name %s
	fb = regexp.MustCompile(`TFMTRESNAME_([sdtfg])`).ReplaceAllString(fb, `%${1}`)

	// resource name %[n]q
	fb = regexp.MustCompile(`"TFMTRESNAME_(\d+)q"`).ReplaceAllString(fb, `%[${1}]q`)

	// resource name %q
	fb = regexp.MustCompile(`"TFMTRESNAME_q"`).ReplaceAllLiteralString(fb, `%q`)

	// count meta-argument
	fb = regexp.MustCompile(`1\s+# tfmtcount_(\[\d+\][ds])`).ReplaceAllString(fb, `%${1}`)
	fb = regexp.MustCompile(`1\s+# tfmtcount_([ds])`).ReplaceAllString(fb, `%${1}`)

	// provider meta-argument
	fb = regexp.MustCompile(`tfmtprovider.PROVIDER_(\d+)`).ReplaceAllString(fb, `%[${1}]s`)
	fb = strings.ReplaceAll(fb, "tfmtprovider.PROVIDER", "%s")

	return fb
}
