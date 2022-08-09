package cli

import (
	"strings"
	"unicode"
)

func indentToOriginalLevel(formatted string, original string) string {
	prefix := ""
	for _, r := range original {
		if unicode.IsSpace(r) {
			if r == '\n' {
				prefix = ""
				continue
			}
			prefix += string(r)
		} else {
			break
		}
	}
	res := strings.ReplaceAll(formatted, "\n", "\n"+prefix)
	res = strings.TrimRight(res, prefix)
	res = strings.TrimLeft(res, prefix)

	return prefix + res
}
