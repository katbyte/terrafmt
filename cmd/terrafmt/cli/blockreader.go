package cli

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/katbyte/terrafmt/common"
	"github.com/sirupsen/logrus"
)

type BlockReader struct {
	FileName string

	//io
	Reader io.Reader
	Writer io.Writer

	//stats
	LineCount  int
	LinesBlock int
	BlockCount int

	ReadOnly bool

	//current block line count
	//blocks formatted
	//

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
	} else if strings.HasPrefix(line, "```hcl") { // documentation
		return true
	}

	return false
}

func IsBlockFinished(line string) bool {
	if line == "`)" { // acctest
		return true
	} else if strings.HasPrefix(line, "`, ") { // acctest
		return true
	} else if strings.HasPrefix(line, "```") { // documentation
		return true
	}

	return false
}

func (br *BlockReader) DoTheThing(filename string) error {

	var tmpfile *os.File

	if filename != "" {
		br.FileName = filename
		common.Log.Debugf("opening src file %s", filename)
		fs, err := os.Open(filename) // For read access.
		if err != nil {
			return err
		}
		defer fs.Close()
		br.Reader = fs

		// for now write to a temporary file, TODO is there a better way?
		if !br.ReadOnly {
			tmpfile, err = ioutil.TempFile("", "terrafmt")
			if err != nil {
				return fmt.Errorf("unable to create tmpfile: %v", err)
			}
			common.Log.Debugf("opening tmp file %s", tmpfile.Name())
			br.Writer = tmpfile
		}

	} else {
		br.FileName = "stdin"
		br.Reader = os.Stdin
		br.Writer = os.Stdout
	}

	br.LineCount = 0
	br.BlockCount = 0
	s := bufio.NewScanner(br.Reader)
	for s.Scan() { //move this to a LineRead function?
		br.LineCount += 1
		//br.CurrentLine = s.Text()+"\n"
		l := s.Text() + "\n"

		if err := br.LineRead(br, br.LineCount, l); err != nil {
			return fmt.Errorf("NB LineRead failed @ %s#%d for %s: %v", br.FileName, br.LineCount, l, err)
		}

		if IsBlockStart(l) {
			block := ""
			br.BlockCount += 1

			for s.Scan() {
				br.LineCount += 1
				l2 := s.Text() + "\n"

				if IsBlockFinished(l2) {
					// todo configure this behaviour with switchs
					if err := br.BlockRead(br, br.LineCount, block); err != nil {
						//for now ignore block errors and output unformatted
						logrus.Errorf("block %d @ %s#%d failed to process with: %v", br.BlockCount, br.FileName, br.LineCount, err)
						BlockReaderPassthrough(br, br.LineCount, block)
					}

					if err := br.LineRead(br, br.LineCount, l2); err != nil {
						return fmt.Errorf("NB LineRead failed @ %s#%d for %s: %v", br.FileName, br.LineCount, l2, err)
					}

					block = ""
					break
				} else {
					block += l2
				}
			}

			/*if block != "" { //couldn't find the end of the block, output as lines
						for each line {
								if err := br.LineRead(br, br.LineCount, l); err != nil {
								return fmt.Errorf("NB LineRead failed @ %s#%d for %s: %v", br.FileName, br.LineCount, l, err)
							}
							}
			 			}*/
		}
	}

	//todo add better error checking and cleanup
	if tmpfile != nil {
		common.Log.Debugf("tmp file %s exists", tmpfile.Name())
		tmpfile.Close()

		common.Log.Debugf("reopening tmp file %s", tmpfile.Name())
		source, err := os.Open(tmpfile.Name())
		if err != nil {
			return err
		}
		defer source.Close()

		common.Log.Debugf("creating destination @ %s", tmpfile.Name())
		destination, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer destination.Close()

		common.Log.Debugf("copying..")
		_, err = io.Copy(destination, source)
		return err
	}

	// todo should this be at the end of a command?
	//fmt.Fprintf(os.Stderr, c.Sprintf("\nFinished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LineCount, br.BlockCount))
	return nil
}
