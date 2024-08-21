package xargs

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Xargs struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// The number of workers to execute commands in parallel with.
	NumWorkers int

	// Whether or not to split input on null-separator (e.g., for usage with find -print0)
	Null bool
}

func (x *Xargs) Run(ctx context.Context, command string, args ...string) {
	wg := sync.WaitGroup{}
	tokens := make(chan string, x.NumWorkers)
	ctx, cancel := context.WithCancel(ctx)
	scanner := bufio.NewScanner(x.Stdin)
	split := bufio.ScanWords
	if x.Null {
		split = splitAtNull
	}
	scanner.Split(split)

	// create workers to run commands
	go func() {
		for range x.NumWorkers {
			go func() {
				for token := range tokens {
					combinedArgs := append(args, token)

					cmd := exec.Command(command, combinedArgs...)
					cmd.Stdin = x.Stdin
					cmd.Stdout = x.Stdout
					cmd.Stderr = x.Stderr

					if err := cmd.Run(); err != nil {
						fmt.Fprintf(x.Stderr, "error running command: %v\n", err)
					}

					wg.Done()
				}

			}()
		}
	}()

	// tokenize input
	for scanner.Scan() {
		wg.Add(1)
		tokens <- scanner.Text()
	}

	go func() {
		wg.Wait()
		close(tokens)
		cancel()
	}()

	<-ctx.Done()
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
