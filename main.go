package main

import (
	"fmt"
	"os"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/cli"
)

func main() {
	if err := cli.Make().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, c.Sprintf("<red>terrafmt:</> %v", err))
		os.Exit(1)
	}

	os.Exit(0)
}
