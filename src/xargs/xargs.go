package xargs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type Xargs struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Whether or not to split input on null-separator (e.g., for usage with find -print0)
	Null bool
}

func (x *Xargs) Run(command string, args ...string) {
	scanner := bufio.NewScanner(x.Stdin)
	split := bufio.ScanWords
	if x.Null {
		split = splitAtNull
	}
	scanner.Split(split)

	for scanner.Scan() {
		combinedArgs := append(args, scanner.Text())

		cmd := exec.Command(command, combinedArgs...)
		cmd.Stdin = x.Stdin
		cmd.Stdout = x.Stdout
		cmd.Stderr = x.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(x.Stderr, err)
		}
	}

	if scanner.Err() != nil {
		fmt.Fprintln(x.Stderr, scanner.Err())
	}
}

// splitAtNull is a bufio.SplitFunc that splits at null characters
func splitAtNull(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}
