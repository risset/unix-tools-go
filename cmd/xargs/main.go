package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	flag.Parse()

	cliArgs := flag.Args()

	if len(cliArgs) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	xargs := &xargs.Xargs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Null:   *null,
	}

	command, args := cliArgs[0], cliArgs[1:]
	xargs.Run(command, args...)
}

func handleInvalidArg(msg string) {
	flag.Usage()
	log.Fatal(msg)
}
