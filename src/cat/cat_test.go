package cat

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	for name, tc := range map[string]struct {
		Paths          []string
		Stdin          string
		ExpectedStdout string
		ExpectedStderr string
	}{
		"stdin only": {
			Stdin:          "a b c",
			ExpectedStdout: "a b c",
		},

		"file only": {
			Paths:          []string{"./testdata/foo.txt"},
			ExpectedStdout: "foo\n",
		},

		"file and stdin": {
			Stdin:          "bar",
			Paths:          []string{"./testdata/foo.txt", "-"},
			ExpectedStdout: "foo\nbar",
		},

		"missing file": {
			Paths:          []string{"missing.txt"},
			ExpectedStderr: "open missing.txt: no such file or directory\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			stdinBuf := bytes.NewBufferString(tc.Stdin)
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			cat := Cat{
				Stdin:  stdinBuf,
				Stdout: &stdoutBuf,
				Stderr: &stderrBuf,
			}

			cat.Run(tc.Paths)

			assert.Equal(t, tc.ExpectedStdout, stdoutBuf.String())
			assert.Equal(t, tc.ExpectedStderr, stderrBuf.String())
		})
	}

}
