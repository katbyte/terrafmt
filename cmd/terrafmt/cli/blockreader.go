package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	c "github.com/gookit/color"
)

type BlockReader struct {
	//io
	Reader io.Reader
	Writer io.Writer

	//stats
	LineCount  int
	LinesBlock int
	BlockCount int

	//callbacks
	LineRead  func(*BlockReader, int, string) error
	BlockRead func(*BlockReader, int, string) error
}

func BlockReaderPassthrough(br *BlockReader, number int, line string) error {
	br.Writer.Write([]byte(line))
	return nil
}

func BlockReaderIgnore(br *BlockReader, number int, line string) error {
	return nil
}

func IsBlockStart(line string) bool {
	if strings.HasSuffix(line, "return fmt.Sprintf(`\n") { // acctest
		return true
	} else if line == "```hcl" { // documentation
		return true
	}

	return false
}

func IsBlockFinished(line string) bool {
	if line == "`)" { // acctest
		return true
	} else if strings.HasPrefix(line, "`, ") { // acctest
		return true
	} else if line == "```" { // documentation
		return true
	}

	return false
}

func (br *BlockReader) DoTheThing() error {
	s := bufio.NewScanner(br.Reader)

	//if writeres are nil default to stdin/stdout?

	br.LineCount = 0
	br.BlockCount = 0
	for s.Scan() { //move this to a LineRead function?
		br.LineCount += 1
		//br.CurrentLine = s.Text()+"\n"
		l := s.Text() + "\n"

		br.LineRead(br, br.LineCount, l)

		if IsBlockStart(l) {
			block := ""
			br.BlockCount += 1

			for s.Scan() {
				br.LineCount += 1
				l2 := s.Text() + "\n"

				if IsBlockFinished(l2) {
					br.BlockRead(br, br.LineCount, block)
					br.LineRead(br, br.LineCount, l2)
					break
				} else {
					block += l2
				}
			}
		}
	}

	// todo should this be at the end of a command?
	fmt.Fprintf(os.Stderr, c.Sprintf("\nFinished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LineCount, br.BlockCount))
	return nil
}
