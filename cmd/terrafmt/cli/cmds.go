package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/common"
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

// reader -> stream -> blocks & non blocks

//flag: comment out %s
// blah = %s^ -> = "$$%s$$"

//reader: start stop pairs
//blocks: ignore %s, ignore ... (docs)

func FormatBlock(b string) (string, error) {

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	common.Log.Debugf("running terraform... ")
	cmd := exec.Command("terraform", "fmt", "-")
	cmd.Stdin = strings.NewReader(b)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed with %s: %s", err, stderr)
	}

	ec := cmd.ProcessState.ExitCode()
	common.Log.Debugf("terraform exited with %d", ec)
	if ec != 0 {
		return "", fmt.Errorf("trerraform failed with %d: %s", ec, stderr)
	}

	return stdout.String(), nil
}

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

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt  %s", filename)

			br := BlockReader{
				LineRead: BlockReaderPassthrough,
				BlockRead: func(br *BlockReader, i int, l string) error {
					fb, err := FormatBlock(l)
					if err != nil {
						return err
					}
					br.Writer.Write([]byte(fb))
					return nil
				},
			}

			err := br.DoTheThing(filename)

			if err != nil {
				return err
			}
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

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt fmt %s", filename)

			br := BlockReader{
				Reader: os.Stdin,
				Writer: os.Stdout,

				LineRead: BlockReaderPassthrough,
				BlockRead: func(br *BlockReader, i int, l string) error {
					fmt.Fprintf(os.Stdout, c.Sprintf("\n<white>#######</> <cyan>B%d</><darkGray> @ #%d</>\n", br.BlockCount, br.LineCount))
					br.Writer.Write([]byte(l))
					return nil
				},
			}

			err := br.DoTheThing(filename)

			if err != nil {
				return err
			}
			return nil
		},
	}
	root.AddCommand(diff)

	// options
	blocks := &cobra.Command{
		Use:   "blocks [file]",
		Short: "extract terraform blocks from a file ",
		Long:  `TODO`,
		//options: no header (######), format (json? xml? ect), only should block x?
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt blocks %s", filename)

			br := BlockReader{
				Reader: os.Stdin,
				Writer: os.Stdout,

				LineRead: BlockReaderIgnore,
				BlockRead: func(br *BlockReader, i int, l string) error {
					fmt.Fprintf(os.Stdout, c.Sprintf("\n<white>%s</>@<white>%d</> <cyan>B%d</>\n", filename, br.LineCount, br.BlockCount))
					br.Writer.Write([]byte(l))
					return nil
				},
			}

			err := br.DoTheThing(filename)

			if err != nil {
				return err
			}

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
