package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andreyvit/diff"
	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/common"
	"github.com/katbyte/terrafmt/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

// reader -> stream -> blocks & non blocks

//flag: comment out %s
// blah = %s^ -> = "$$%s$$"

//reader: start stop pairs
//blocks: ignore %s, ignore ... (docs)

//stats: lines, blocks, blocks formatted (lines formatted?), errors?
func Make() *cobra.Command {

	//shared options
	// acctest mode (format string compatible)
	// quite
	// verbose?

	root := &cobra.Command{
		Use:   "terrafmt [fmt|diff|blocks]",
		Short: "terrafmt is a small utility to format terraform blocks found in files.",
		Long: `A small utility to for formatting terraform blocks found in files. Primarily intended to help with terraform provider development.
Complete documentation is available at https://github.com/katbyte/terrafmt`,
		Args:          cobra.RangeArgs(0, 0),
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			return fmt.Errorf("No command specified")
		},
	}

	//options : only count, blocks diff/found, total lines diff, etc
	root.AddCommand(&cobra.Command{
		Use:   "fmt [file]",
		Short: "formats terraform blocks in a single file or on stdin",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt  %s", filename)

			blocksFormatted := 0
			br := BlockReader{
				LineRead: BlockReaderPassthrough,
				BlockRead: func(br *BlockReader, i int, b string) error {
					fb, err := FormatBlock(b, viper.GetBool("fmtcompat"))
					if err != nil {
						return err
					}
					br.Writer.Write([]byte(fb))

					if fb != b {
						blocksFormatted++
					}

					return nil
				},
			}
			err := br.DoTheThing(filename)

			fmt.Fprintf(os.Stderr, c.Sprintf("Processed <magenta>%s</>: <cyan>%d</> lines & formatted <yellow>%d</>/<yellow>%d</> blocks!\n", br.FileName, br.LineCount, blocksFormatted, br.BlockCount))

			if err != nil {
				return err
			}
			return nil
		},
	})

	//options : only count, blocks diff/found, total lines diff, etc
	root.AddCommand(&cobra.Command{
		Use:   "diff [file]",
		Short: "formats terraform blocks in a file and shows the difference",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt fmt %s", filename)

			blocksWithDiff := 0
			br := BlockReader{
				ReadOnly: true,
				LineRead: func(br *BlockReader, i int, b string) error {
					return nil
				},
				BlockRead: func(br *BlockReader, i int, b string) error {
					fb, err := FormatBlock(b, viper.GetBool("fmtcompat"))
					if err != nil {
						return err
					}

					if fb == b {
						return nil
					}
					blocksWithDiff++

					fmt.Fprintf(os.Stderr, c.Sprintf("\n<white>#######</> <cyan>B%d</><darkGray> @</> <lightMagenta>%s</><darkGray>#</><magenta>%d</>\n", br.BlockCount, br.FileName, br.LineCount-br.BlockCurrentLine))

					d := diff.LineDiff(b, fb)
					scanner := bufio.NewScanner(strings.NewReader(d))
					for scanner.Scan() {
						l := scanner.Text()
						if strings.HasPrefix(l, "+") {
							fmt.Fprint(os.Stdout, c.Sprintf("<green>%s</>\n", l))
						} else if strings.HasPrefix(l, "-") {
							fmt.Fprint(os.Stdout, c.Sprintf("<red>%s</>\n", l))
						} else {
							fmt.Fprint(os.Stdout, l+"\n")
						}
					}

					return nil
				},
			}
			err := br.DoTheThing(filename)

			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, c.Sprintf("\nProcessed <magenta>%s</>: <cyan>%d</> lines & <yellow>%d</>/<yellow>%d</> blocks need formatting.\n", br.FileName, br.LineCount, blocksWithDiff, br.BlockCount))

			return nil
		},
	})

	// options
	root.AddCommand(&cobra.Command{
		Use:   "blocks [file]",
		Short: "extracts terraform blocks from a file ",
		//options: no header (######), format (json? xml? ect), only should block x?
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt blocks %s", filename)

			br := BlockReader{
				ReadOnly: true,
				LineRead: BlockReaderIgnore,
				BlockRead: func(br *BlockReader, i int, b string) error {
					fmt.Fprintf(os.Stdout, c.Sprintf("\n<white>#######</> <cyan>B%d</><darkGray> @ #%d</>\n", br.BlockCount, br.LineCount))
					fmt.Fprint(os.Stdout, b)
					return nil
				},
			}

			err := br.DoTheThing(filename)

			if err != nil {
				return err
			}

			//blocks
			fmt.Fprintf(os.Stderr, c.Sprintf("\nFinished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LineCount, br.BlockCount))

			return nil
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of terrafmt",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("terrafmt v" + version.Version + "-" + version.GitCommit)
		},
	})

	pflags := root.PersistentFlags()
	pflags.BoolP("fmtcompat", "f", false, "enable format string (%s, %d ect) compatibility")

	viper.BindPFlag("fmtcompat", pflags.Lookup("fmtcompat"))

	//todo bind to env?

	return root
}
