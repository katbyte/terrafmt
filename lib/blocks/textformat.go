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
	//nolint:gocritic
	if strings.HasPrefix(line, "```hcl") { // documentation
		return true
	} else if strings.HasPrefix(line, "```terraform") { // documentation
		return true
	} else if strings.HasPrefix(line, "```tf") { // documentation
		return true
	}

	return false
}

func (mbf markdownTextFormat) isFinishLine(line string) bool {
	return strings.HasPrefix(line, "```")
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
	return strings.Compare(line, strings.TrimLeftFunc(line, unicode.IsSpace)) == 0
}

func (mbf restructuredTextFormat) preserveIndentation() bool {
	return true
}
