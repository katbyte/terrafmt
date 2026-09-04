package cli

import (
	"strings"
	"unicode"
)

func indentToOriginalLevel(formatted, original string) string {
	prefix := ""
	for _, r := range original {
		if !unicode.IsSpace(r) {
			break
		}
		if r == '\n' {
			prefix = ""
			continue
		}
		prefix += string(r)
	}
	res := strings.ReplaceAll(formatted, "\n", "\n"+prefix)
	res = strings.TrimRight(res, prefix)
	res = strings.TrimLeft(res, prefix)

	return prefix + res
}
