package main

import (
	"fmt"
	"os"

	c "github.com/gookit/color"
	"github.com/katbyte/terrafmt/cli"
)

func main() {
	root, err := cli.Make()
	if err != nil {
		fmt.Fprint(os.Stderr, c.Sprintf("<red>terrafmt:</> %v\n", err))
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Fprint(os.Stderr, c.Sprintf("<red>terrafmt:</> %v\n", err))
		os.Exit(1)
	}

	os.Exit(0)
}
