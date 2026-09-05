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

	// documentation fences
	return strings.HasPrefix(trimmed, "```hcl") ||
		strings.HasPrefix(trimmed, "```terraform") ||
		strings.HasPrefix(trimmed, "```tf")
}

func (markdownTextFormat) isFinishLine(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "```")
}

func (markdownTextFormat) preserveIndentation() bool {
	return false
}

// used for restructured text
type restructuredTextFormat struct{}

func (restructuredTextFormat) isStartingLine(line string) bool {
	return strings.HasPrefix(line, ".. code:: terraform")
}

func (restructuredTextFormat) isFinishLine(line string) bool {
	return line == strings.TrimLeftFunc(line, unicode.IsSpace)
}

func (restructuredTextFormat) preserveIndentation() bool {
	return true
}
