package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andreyvit/diff"
	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/lib/blocks"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/katbyte/terrafmt/lib/format"
	"github.com/katbyte/terrafmt/lib/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Make() *cobra.Command {
	root := &cobra.Command{
		Use:           "terrafmt [fmt|diff|blocks]",
		Short:         "terrafmt is a small utility to format terraform blocks found in files.",
		Long:          `A small utility that formats terraform blocks found in files. Primarily intended to help with terraform provider development.`,
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
			br := blocks.Reader{
				LineRead: blocks.ReaderPassthrough,
				BlockRead: func(br *blocks.Reader, i int, b string) error {
					var fb string
					var err error
					if viper.GetBool("fmtcompat") {
						fb, err = format.FmtVerbBlock(b)
					} else {
						fb, err = format.Block(b)
					}

					if err != nil {
						return err
					}

					_, err = br.Writer.Write([]byte(fb))

					if err == nil && fb != b {
						blocksFormatted++
					}

					return err
				},
			}
			err := br.DoTheThing(filename)

			fc := "magenta"
			if blocksFormatted > 0 {
				fc = "lightMagenta"
			}

			if !viper.GetBool("quiet") {
				// nolint staticcheck
				fmt.Fprintf(os.Stderr, c.Sprintf("<%s>%s</>: <cyan>%d</> lines & formatted <yellow>%d</>/<yellow>%d</> blocks!\n", fc, br.FileName, br.LineCount, blocksFormatted, br.BlockCount))
			}
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
			br := blocks.Reader{
				ReadOnly: true,
				LineRead: blocks.ReaderPassthrough,
				BlockRead: func(br *blocks.Reader, i int, b string) error {
					var fb string
					var err error
					if viper.GetBool("fmtcompat") {
						fb, err = format.FmtVerbBlock(b)
					} else {
						fb, err = format.Block(b)
					}

					if err != nil {
						return err
					}

					if fb == b {
						return nil
					}
					blocksWithDiff++

					// nolint staticcheck
					fmt.Fprintf(os.Stdout, c.Sprintf("<lightMagenta>%s</><darkGray>#</><magenta>%d</>\n", br.FileName, br.LineCount-br.BlockCurrentLine))

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

			fc := "magenta"
			if blocksWithDiff > 0 {
				fc = "lightMagenta"
			}

			if !viper.GetBool("quiet") {
				// nolint staticcheck
				fmt.Fprintf(os.Stderr, c.Sprintf("<%s>%s</>: <cyan>%d</> lines & <yellow>%d</>/<yellow>%d</> blocks need formatting.\n", fc, br.FileName, br.LineCount, blocksWithDiff, br.BlockCount))
			}
			return nil
		},
	})

	// options
	root.AddCommand(&cobra.Command{
		Use:   "blocks [file]",
		Short: "extracts terraform blocks from a file ",
		//options: no header (######), format (json? xml? etc), only should block x?
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}
			common.Log.Debugf("terrafmt blocks %s", filename)

			br := blocks.Reader{
				ReadOnly: true,
				LineRead: blocks.ReaderIgnore,
				BlockRead: func(br *blocks.Reader, i int, b string) error {
					// nolint staticcheck
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
			// nolint staticcheck
			fmt.Fprintf(os.Stderr, c.Sprintf("\nFinished processing <cyan>%d</> lines <yellow>%d</> blocks!\n", br.LineCount, br.BlockCount))

			return nil
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of terrafmt",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			//nolint errcheck
			fmt.Println("terrafmt v" + version.Version + "-" + version.GitCommit)
		},
	})

	pflags := root.PersistentFlags()
	pflags.BoolP("fmtcompat", "f", false, "enable format string (%s, %d etc) compatibility")
	pflags.BoolP("quiet", "q", false, "only show differences")

	if err := viper.BindPFlag("fmtcompat", pflags.Lookup("fmtcompat")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("quiet", pflags.Lookup("quiet")); err != nil {
		panic(err)
	}

	//todo bind to env?

	return root
}
