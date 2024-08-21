package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/risset/unix-utils/src/find"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	null := flag.Bool("0", false, "separate files by null character instead of newline")
	flag.Parse()
	cliArgs := flag.Args()

	var dir, pattern string
	switch len(cliArgs) {
	case 1:
		dir = cliArgs[0]

	case 2:
		dir, pattern = cliArgs[0], cliArgs[1]

	default:
		fmt.Fprintln(os.Stderr, "Usage: find <dir> <pattern>")
		os.Exit(1)
	}

	separator := "\n"
	if *null {
		separator = "\x00"
	}

	find := &find.Find{
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Dir:       dir,
		Pattern:   pattern,
		Separator: separator,
	}

	if err := find.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
