package xargs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	for name, tc := range map[string]struct {
		Stdin          string
		ExpectedStdout string
		ExpectedStderr string
		Null           bool
		Command        string
		Args           []string
	}{
		"echo with no args": {
			Command:        "echo",
			Args:           []string{},
			Stdin:          "foo bar",
			ExpectedStdout: "foo\nbar\n",
		},

		"echo with args": {
			Command:        "echo",
			Args:           []string{"baz"},
			Stdin:          "foo bar",
			ExpectedStdout: "baz foo\nbaz bar\n",
		},

		"handle null-separated input": {
			Command:        "echo",
			Args:           []string{"baz"},
			Stdin:          "foo\x00bar",
			Null:           true,
			ExpectedStdout: "baz foo\nbaz bar\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			stdinBuf := bytes.NewBufferString(tc.Stdin)
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			x := Xargs{
				Stdin:  stdinBuf,
				Stdout: &stdoutBuf,
				Stderr: &stderrBuf,
				Null:   tc.Null,
			}

			x.Run(tc.Command, tc.Args...)

			assert.Equal(t, tc.ExpectedStdout, stdoutBuf.String())
			assert.Equal(t, tc.ExpectedStderr, stderrBuf.String())
		})
	}
}
