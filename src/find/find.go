package find

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Find struct {
	Stdout io.Writer
	Stderr io.Writer

	// The directory to find files in
	Dir string

	// The regex pattern to match files against (optional)
	Pattern string

	// The separator to separate file results with (defaults to "\n")
	Separator string
}

// Run executes the find command.
func (f *Find) Run(ctx context.Context) error {
	if f.Separator == "" {
		f.Separator = "\n"
	}

	paths, errs := walk(f.Dir)
	if f.Pattern != "" {
		re, err := regexp.Compile(f.Pattern)
		if err != nil {
			return err
		}

		paths = filter(paths, re)
	}

	for {
		select {
		case path, ok := <-paths:
			if !ok {
				return nil
			}

			fmt.Fprint(f.Stdout, path+f.Separator)

		case err := <-errs:
			if os.IsNotExist(err) {
				fmt.Fprintf(f.Stderr, "%q does not exist\n", f.Dir)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// walk does a recursive walk over the given directory
func walk(dir string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errs := make(chan error)

	go func() {
		defer close(paths)

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if dir != "" && dir[:2] == "./" && path != dir {
				path = "./" + path
			}

			if err == nil {
				paths <- path
			}

			return err
		})

		if err != nil {
			errs <- err
		}
	}()

	return paths, errs
}

// filter filters the given paths based on a regular expression
func filter(paths <-chan string, re *regexp.Regexp) <-chan string {
	filteredPaths := make(chan string)

	go func() {
		defer close(filteredPaths)

		for p := range paths {
			if re.MatchString(p) {
				filteredPaths <- p
			}
		}
	}()

	return filteredPaths
}
