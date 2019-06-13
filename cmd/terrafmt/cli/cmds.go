package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ValidateParams(params []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, p := range params {
			if viper.GetString(p) == "" {
				return fmt.Errorf(p + " paramter can't be empty")
			}
		}

		return nil
	}
}

// data -> read -> chunk block & non block
//  non block -> passthrough
//  block     -> act on blocks
// combine -> final stream

// reader:
// non block line: blah
// block line ?
// block start ?
// block finished
//

type BlockReader struct {
	//io
	Reader io.Reader
	Writer io.Writer

	//stats
	LinesCount int
	LinesBlock int
	BlockCount int

	//callbacks
	ReadLine  func(*BlockReader, int, string) error
	ReadBlock func(*BlockReader, int, string) error
}

func BlockReaderPassthrough(br *BlockReader, number int, line string) error {
	br.Writer.Write([]byte(line))
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
	s := bufio.NewScanner(os.Stdin)

	br.LinesCount = 0
	br.BlockCount = 0
	for s.Scan() { //move this to a ReadLine function?
		br.LinesCount += 1
		//br.CurrentLine = s.Text()+"\n"
		l := s.Text() + "\n"

		br.ReadLine(br, br.LinesCount, l)

		if IsBlockStart(l) {
			block := ""
			br.BlockCount += 1

			for s.Scan() {
				br.LinesCount += 1
				l2 := s.Text() + "\n"

				if IsBlockFinished(l2) {
					br.ReadLine(br, br.LinesCount, l2)
					break
				} else {
					block += l2
					br.ReadBlock(br, br.LinesCount, block)
					fmt.Fprint(os.Stderr, block)
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, c.Sprintf("Finished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LinesCount, br.BlockCount))
	return nil
}

// reader -> stream -> blocks & non blocks

//flag: comment out %s
// blah = %s^ -> = "$$%s$$"

//reader: start stop pairs
//blocks: ignore %s, ignore ... (docs)

//stats: lines, blocks, blocks formatted (lines formatted?), errors?
func Make() *cobra.Command {

	root := &cobra.Command{
		Use:   "terrafmt [file]",
		Short: "terrafmt is a small utility to trigger acceptance tests on teamcity",
		Long: `A small utility to trigger acceptance tests on teamcity. 
It can also pull the tests to run for a PR on github
Complete documentation is available at https://github.com/katbyte/terrafmt`,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("terrafmt bare")

			br := BlockReader{
				Reader: os.Stdin,
				Writer: os.Stdout,

				ReadLine:  BlockReaderPassthrough,
				ReadBlock: BlockReaderPassthrough,
			}
			err := br.DoTheThing()

			if err != nil {
				return err
			}

			//reader to read file and find blocks

			//fmt blocks

			//output file + blocks

			return nil
		},
	}

	//options : only count, blocks diff/found, total lines diff, etc
	diff := &cobra.Command{
		Use:   "diff [file]",
		Short: "formats terraform blocks in a file and shows the differnce",
		Long:  `TODO`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("terrafmt diff")

			//reader to read file and find blocks

			//fmt blocks

			//blocks that are different etc
			return nil
		},
	}
	root.AddCommand(diff)

	// options
	blocks := &cobra.Command{
		Use:   "blocks [file]",
		Short: "extract terraform blocks from a file ",
		Long:  `TODO`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("terrafmt block")
			//reader to read file and find blocks

			//blocks
			return nil
		},
	}
	root.AddCommand(blocks)

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of terrafmt",
		Long:  `Print the version number of terrafmt`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("terrafmt v" + version.Version + "-" + version.GitCommit)
		},
	})

	return root
}
