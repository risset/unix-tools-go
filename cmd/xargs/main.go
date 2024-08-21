package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/risset/unix-utils/src/xargs"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s <command> <args...>:\n", os.Args[0])
		flag.PrintDefaults()
	}

	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	null := flag.Bool("0", false, "delimit by null character instead of whitespace")
	numWorkers := flag.Int("w", 10, "number of workers to execute commands in parallel")

	flag.Parse()

	if *numWorkers < 1 {
		handleInvalidArg("number of workers must be >= 1")
	}

	cliArgs := flag.Args()

	if len(cliArgs) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	xargs := &xargs.Xargs{
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		NumWorkers: *numWorkers,
		Null:       *null,
	}

	command, args := cliArgs[0], cliArgs[1:]
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	xargs.Run(ctx, command, args...)
}

func handleInvalidArg(msg string) {
	flag.Usage()
	log.Fatal(msg)
}
