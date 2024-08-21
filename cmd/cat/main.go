package main

import (
	"os"

	"github.com/risset/unix-utils/src/cat"
)

func main() {
	cat := &cat.Cat{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var paths []string
	switch len(os.Args) {
	case 0, 1:
	default:
		paths = os.Args[1:]
	}

	cat.Run(paths)
}
