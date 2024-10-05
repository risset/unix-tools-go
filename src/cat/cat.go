package cat

import (
	"fmt"
	"io"
	"os"
)

type Cat struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Cat) Run(files []string) {
	if len(files) == 0 {
		_, err := io.Copy(c.Stdout, c.Stdin)
		if err != nil {
			fmt.Fprintln(c.Stderr, err)
		}
	}

	for _, path := range files {
		switch path {
		case "-":
			_, err := io.Copy(c.Stdout, c.Stdin)
			if err != nil {
				fmt.Fprintln(c.Stderr, err)
			}

		default:
			file, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(c.Stderr, err)
				continue
			}

			_, err = io.Copy(c.Stdout, file)
			if err != nil {
				fmt.Fprintln(c.Stderr, err)
			}
		}
	}
}
