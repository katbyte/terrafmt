package blocks

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	accTestFinishLineWithLeadingSpacesMatcher = regexp.MustCompile("^[[:space:]]*`(,|\\)\n)")
	lineWithLeadingSpacesMatcher              = regexp.MustCompile("^[[:space:]]*(.*\n)$")
)

type BlockWriter interface {
	Write(index, startLine, endLine int, text string)
	Close() error
}

type Reader struct {
	FileName string

	Log *logrus.Logger

	//io
	Reader io.Reader
	Writer io.Writer

	//stats
	LineCount        int // total lines processed
	LinesBlock       int // total block lines processed
	BlockCount       int // total blocks found
	BlockCurrentLine int // current block line count

	ErrorBlocks int

	//options
	ReadOnly       bool
	FixFinishLines bool

	//callbacks
	LineRead  func(*Reader, int, string) error
	BlockRead func(*Reader, int, string) error

	// Only used by the "blocks" command
	BlockWriter BlockWriter
}

func ReaderPassthrough(br *Reader, number int, line string) error {
	_, err := br.Writer.Write([]byte(line))
	return err
}

func ReaderIgnore(br *Reader, number int, line string) error {
	return nil
}

func IsStartLine(line string) bool {
	if strings.HasSuffix(line, "return fmt.Sprintf(`\n") { // acctest
		return true
	} else if strings.HasPrefix(line, "```hcl") { // documentation
		return true
	} else if strings.HasPrefix(line, "```terraform") { // documentation
		return true
	} else if strings.HasPrefix(line, "```tf") { // documentation
		return true
	}

	return false
}

func IsFinishLine(line string) bool {
	if accTestFinishLineWithLeadingSpacesMatcher.MatchString(line) { // acctest
		return true
	} else if strings.HasPrefix(line, "```") { // documentation
		return true
	}

	return false
}

func (br *Reader) DoTheThing(fs afero.Fs, filename string, stdin io.Reader, stdout io.Writer) error {
	var buf *bytes.Buffer

	if filename != "" {
		br.FileName = filename
		br.Log.Debugf("opening src file %s", filename)
		file, err := fs.Open(filename) // For read access.
		if err != nil {
			return err
		}
		defer file.Close()
		br.Reader = file

		// for now write to buffer
		if !br.ReadOnly {
			buf = bytes.NewBuffer([]byte{})
			br.Writer = buf
		} else {
			br.Writer = ioutil.Discard
		}
	} else {
		br.FileName = "stdin"
		br.Reader = stdin
		br.Writer = stdout

		if br.ReadOnly {
			br.Writer = ioutil.Discard
		}
	}

	br.LineCount = 0
	br.BlockCount = 0
	s := bufio.NewScanner(br.Reader)
	for s.Scan() { // scan file
		br.LineCount += 1
		//br.CurrentLine = s.Text()+"\n"
		l := s.Text() + "\n"

		if err := br.LineRead(br, br.LineCount, l); err != nil {
			return fmt.Errorf("NB LineRead failed @ %s:%d for %s: %v", br.FileName, br.LineCount, l, err)
		}

		if IsStartLine(l) {
			block := ""
			br.BlockCurrentLine = 0
			br.BlockCount += 1

			for s.Scan() { // scan block
				br.LineCount += 1
				br.BlockCurrentLine += 1
				l2 := s.Text() + "\n"

				// make sure we don't run into another block
				if IsStartLine(l2) {
					// the end of current block must be malformed, so lets pass it through and log an error
					br.Log.Errorf("block %d @ %s:%d failed to find end of block", br.BlockCount, br.FileName, br.LineCount-br.BlockCurrentLine)
					if err := ReaderPassthrough(br, br.LineCount, block); err != nil { // is this ok or should we loop with LineRead?
						return err
					}

					if err := br.LineRead(br, br.LineCount, l2); err != nil {
						return fmt.Errorf("NB LineRead failed @ %s:%d for %s: %v", br.FileName, br.LineCount, l, err)
					}

					block = ""
					br.BlockCount += 1
					continue
				}

				if IsFinishLine(l2) {
					if br.FixFinishLines {
						l2 = lineWithLeadingSpacesMatcher.ReplaceAllString(l2, `$1`)
					}

					br.LinesBlock += br.BlockCurrentLine

					// todo configure this behaviour with switch's
					if err := br.BlockRead(br, br.LineCount, block); err != nil {
						//for now ignore block errors and output unformatted
						br.ErrorBlocks += 1
						br.Log.Errorf("block %d @ %s:%d failed to process with: %v", br.BlockCount, br.FileName, br.LineCount-br.BlockCurrentLine, err)
						if err := ReaderPassthrough(br, br.LineCount, block); err != nil {
							return err
						}
					}

					if err := br.LineRead(br, br.LineCount, l2); err != nil {
						return fmt.Errorf("NB LineRead failed @ %s:%d for %s: %v", br.FileName, br.LineCount, l2, err)
					}

					block = ""
					break
				} else {
					block += l2
				}
			}

			// ensure last block in the file was property handled
			if block != "" {
				//for each line { Lineread()?
				br.Log.Errorf("block %d @ %s:%d failed to find end of block", br.BlockCount, br.FileName, br.LineCount-br.BlockCurrentLine)
				if err := ReaderPassthrough(br, br.LineCount, block); err != nil { // is this ok or should we loop with LineRead?
					return err
				}
			}
		}
	}

	// If not read-only, need to write back to file.
	if !br.ReadOnly && filename != "" {
		destination, err := fs.Create(filename)
		if err != nil {
			return err
		}
		defer destination.Close()

		br.Log.Debugf("copying..")
		_, err = io.Copy(destination, buf)
		return err
	}

	// todo should this be at the end of a command?
	//fmt.Fprintf(os.Stderr, c.Sprintf("\nFinished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LineCount, br.BlockCount))
	return nil
}
