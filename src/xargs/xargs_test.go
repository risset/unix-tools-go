package xargs

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	for name, tc := range map[string]struct {
		Stdin          string
		ExpectedStdout string
		ExpectedStderr string
		NumWorkers     int
		Null           bool
		Command        string
		Args           []string
	}{
		"echo with no args": {
			Command:        "echo",
			Args:           []string{},
			NumWorkers:     1,
			Stdin:          "foo bar",
			ExpectedStdout: "foo\nbar\n",
		},

		"echo with args": {
			Command:        "echo",
			Args:           []string{"baz"},
			NumWorkers:     1,
			Stdin:          "foo bar",
			ExpectedStdout: "baz foo\nbaz bar\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			stdinBuf := bytes.NewBufferString(tc.Stdin)
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			x := Xargs{
				Stdin:      stdinBuf,
				Stdout:     &stdoutBuf,
				Stderr:     &stderrBuf,
				NumWorkers: tc.NumWorkers,
				Null:       tc.Null,
			}

			ctx := context.Background()
			x.Run(ctx, tc.Command, tc.Args...)

			assert.Equal(t, tc.ExpectedStdout, stdoutBuf.String())
			assert.Equal(t, tc.ExpectedStdout, stdoutBuf.String())
		})
	}

}
