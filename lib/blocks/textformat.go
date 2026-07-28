package blocks

import (
	"strings"
	"unicode"
)

type textFormat interface {
	isStartingLine(line string) bool
	isFinishLine(line string) bool
	preserveIndentation() bool
}

// used for markdown text
type markdownTextFormat struct{}

func (markdownTextFormat) isStartingLine(line string) bool {
	// fences may be indented, e.g. inside a list item (issue #51)
	trimmed := strings.TrimLeft(line, " \t")

	//nolint:gocritic
	if strings.HasPrefix(trimmed, "```hcl") { // documentation
		return true
	} else if strings.HasPrefix(trimmed, "```terraform") { // documentation
		return true
	} else if strings.HasPrefix(trimmed, "```tf") { // documentation
		return true
	}

	return false
}

func (mbf markdownTextFormat) isFinishLine(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "```")
}

func (mbf markdownTextFormat) preserveIndentation() bool {
	return false
}

// used for restructured text
type restructuredTextFormat struct{}

func (restructuredTextFormat) isStartingLine(line string) bool {
	return strings.HasPrefix(line, ".. code:: terraform")
}

func (mbf restructuredTextFormat) isFinishLine(line string) bool {
	return line == strings.TrimLeftFunc(line, unicode.IsSpace)
}

func (mbf restructuredTextFormat) preserveIndentation() bool {
	return true
}
